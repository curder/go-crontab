package worker

import (
    "context"
    "github.com/curder/go-crontab/crontab/common"
    "go.etcd.io/etcd/clientv3"
    "go.etcd.io/etcd/mvcc/mvccpb"
    "time"
)

type JobMgr struct {
    client  *clientv3.Client
    kv      clientv3.KV
    lease   clientv3.Lease
    watcher clientv3.Watcher
}

var (
    GJobMgr *JobMgr // 单例
)

// 初始化etct
func InitJobMgr() (err error) {
    var (
        config  clientv3.Config
        client  *clientv3.Client
        kv      clientv3.KV
        lease   clientv3.Lease
        watcher clientv3.Watcher
    )

    // 初始化配置
    config = clientv3.Config{
        Endpoints:   GConfig.EtcdEndPoints,
        DialTimeout: time.Duration(GConfig.EtcdDialTimeout) * time.Millisecond,
    }

    // 建立连接
    if client, err = clientv3.New(config); err != nil {
        return
    }

    // 得到KV和Lease的API子集
    kv = clientv3.NewKV(client)
    lease = clientv3.NewLease(client)
    watcher = clientv3.NewWatcher(client)

    // 赋值单例
    GJobMgr = &JobMgr{
        client:  client,
        kv:      kv,
        lease:   lease,
        watcher: watcher,
    }

    if err = GJobMgr.watchJobs(); err != nil {
        return
    }

    return
}

// 监听任务变化
func (j *JobMgr) watchJobs() (err error) {
    var (
        getResponse        *clientv3.GetResponse
        kvpair             *mvccpb.KeyValue
        job                *common.Job
        watchStartRevision int64
        watchChan          clientv3.WatchChan
        watchResponse      clientv3.WatchResponse
        watchEvent         clientv3.Event
        jobName            string
        jobEvent           *common.JobEvent
    )
    // 获取/con/jobs/目录下所有任务，并获得当前集群的revision
    if getResponse, err = j.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix()); err != nil {
        return
    }

    // 当前任务
    for _, kvpair = range getResponse.Kvs {
        // 但序列化任务得到job
        if job, err = common.UnpackJob(kvpair.Value); err != nil {
            // 构建任务
            jobEvent = common.BuildJobEvent(common.JobEventSave, job)

            // TODO 把任务同步给Scheduler（调度协程）
        }
    }

    // 从revision向后监听事件
    go func() { // 监听协程
        // 从当前版本的后续版本监听
        watchStartRevision = getResponse.Header.Revision + 1
        // 监听/cron/jobs目录变化
        watchChan = j.watcher.Watch(context.TODO(), common.JobSaveDir, clientv3.WithRev(watchStartRevision))

        // 处理监听事件
        for watchResponse = range watchChan {
            for watchEvent = range watchResponse.Events {
                switch watchEvent.Type {
                case mvccpb.PUT: // 保存任务
                    // 反序列化
                    if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
                        continue // 忽略错误
                    }

                    // 构造一个event事件
                    jobEvent = common.BuildJobEvent(common.JobEventSave, job)

                    // 推送给scheduler
                case mvccpb.DELETE: // 删除任务
                    // todo 推送一个删除事件给scheduler

                    jobName = common.ExtraJobName(string(watchEvent.Kv.Key))

                    job = &common.Job{Name: jobName,} // 构建一个包含名称的Job
                    //  构建删除event事件
                    jobEvent = common.BuildJobEvent(common.JobEventDelete, job)

                    // TODO 推送给scheduler
                }
            }
        }
    }()

    return
}
