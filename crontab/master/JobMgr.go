package master

import (
    "context"
    "encoding/json"
    "github.com/curder/go-crontab/crontab/common"
    "go.etcd.io/etcd/clientv3"
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
    jobKey = `/cron/jobs/` + job.Name

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
