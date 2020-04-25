package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "go.etcd.io/etcd/mvcc/mvccpb"
    "time"
)

const KEY = `/cron/jobs/key1`
const VALUE = `value1`

// 监听KV变化
func main() {
    var (
        config             clientv3.Config
        client             *clientv3.Client
        err                error
        kv                 clientv3.KV
        getResponse        *clientv3.GetResponse
        watchStartRevision int64
        watcher            clientv3.Watcher
        watchChan          clientv3.WatchChan
        watchResponse      clientv3.WatchResponse
        events             *clientv3.Event
    )

    // 客户端连接配置
    config = clientv3.Config{
        Endpoints:   []string{"192.168.205.10:2379"}, // etcd服务地址
        DialTimeout: 5 * time.Second,                 // 建立连接的超时时间
    }

    // 建立连接
    if client, err = clientv3.New(config); err != nil {
        fmt.Printf("Failed to connect etcd: %s", err.Error())
        return
    }

    // 延迟关闭连接
    defer client.Close()

    // KV
    kv = clientv3.NewKV(client)

    // 协程模拟ectd中KV的变化
    go func() {
        for {
            _, _ = kv.Put(context.TODO(), KEY, VALUE) // 写入KV
            _, _ = kv.Delete(context.TODO(), KEY)     // 删除KV

            time.Sleep(1 * time.Second)
        }
    }()

    // 监听key的变化
    // 先GET到对应key当前的值，并监听后续变化
    if getResponse, err = kv.Get(context.TODO(), KEY); err != nil {
        fmt.Printf("Failed to get KV, err: %s\n", err.Error())
        return
    }

    // 存在KV
    if len(getResponse.Kvs) != 0 {
        fmt.Printf("当前值：%s\n", getResponse.Kvs[0].Value)
    }

    // 当前exct集群中的事务ID，递增值
    watchStartRevision = getResponse.Header.Revision + 1

    // 创建一个watcher
    watcher = clientv3.NewWatcher(client)

    // 启动监听
    fmt.Printf("从该版本后监听：%d \n", watchStartRevision)

    watchChan = watcher.Watch(context.TODO(), KEY, clientv3.WithRev(watchStartRevision))

    for watchResponse = range watchChan {
        for _, events = range watchResponse.Events {
            switch events.Type {
            case mvccpb.PUT:
                fmt.Printf("修改为：%s, Revision: %d, Revision: %d\n", events.Kv.Value, events.Kv.CreateRevision, events.Kv.ModRevision)
            case mvccpb.DELETE:
                fmt.Printf("删除KV，Revision: %d\n", events.Kv.ModRevision)
            }
        }
    }
}
