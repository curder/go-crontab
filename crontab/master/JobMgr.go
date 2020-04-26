package master

import (
    "context"
    "encoding/json"
    "github.com/curder/go-crontab/crontab/common"
    "go.etcd.io/etcd/clientv3"
    "go.etcd.io/etcd/mvcc/mvccpb"
    "time"
)

type JobMgr struct {
    client *clientv3.Client
    kv     clientv3.KV
    lease  clientv3.Lease
}

var (
    GJobMgr *JobMgr // 单例
)

// 初始化etct
func InitJobMgr() (err error) {
    var (
        config clientv3.Config
        client *clientv3.Client
        kv     clientv3.KV
        lease  clientv3.Lease
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

    // 赋值单例
    GJobMgr = &JobMgr{
        client: client,
        kv:     kv,
        lease:  lease,
    }

    return
}

// 保存任务
func (j *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
    // 把任务保存到 /cron/jobs/JOB_NAME 下
    var (
        jobKey   string
        jobValue []byte

        putResponse *clientv3.PutResponse

        oldJobObject common.Job
    )

    // etcd中保存的key
    jobKey = common.JOB_SAVE_DIR + job.Name

    // 任务信息
    if jobValue, err = json.Marshal(job); err != nil {
        return
    }

    // 保存到etcd
    if putResponse, err = j.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
        return
    }

    // 如果是更新，返回旧值
    if putResponse.PrevKv != nil {
        if err = json.Unmarshal(putResponse.PrevKv.Value, &oldJobObject); err != nil {
            err = nil
            return
        }
        oldJob = &oldJobObject
    }
    return
}

// 删除任务
func (j *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
    var (
        jobKey         string
        deleteResponse *clientv3.DeleteResponse
        oldJobObject   common.Job
    )

    // etcd中保存的任务key
    jobKey = common.JOB_SAVE_DIR + name

    // 从etcd中删除对应key
    if deleteResponse, err = j.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
        return
    }

    // 返回被删除的任务信息
    if len(deleteResponse.PrevKvs) != 0 {
        // 解析旧的KV，返回它
        if err = json.Unmarshal(deleteResponse.PrevKvs[0].Value, &oldJobObject); err != nil {
            err = nil // 忽略错误
            return
        }
        oldJob = &oldJobObject
    }

    return
}

// 任务列表
func (j *JobMgr) ListJobs() (jobList []*common.Job, err error) {
    var (
        dirKey      string
        getResponse *clientv3.GetResponse
        kvPair      *mvccpb.KeyValue
        job         *common.Job
    )

    // 任务保存的目录
    dirKey = common.JOB_SAVE_DIR

    // 获取目录下的所有任务信息
    if getResponse, err = j.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
        return
    }

    // 初始化数组空间
    jobList = make([]*common.Job, 0)

    // 遍历所有任务并反序列化
    for _, kvPair = range getResponse.Kvs {
        job = &common.Job{}
        if err = json.Unmarshal(kvPair.Value, job); err != nil {
            err = nil
            continue
        }
        jobList = append(jobList, job)
    }

    return
}
