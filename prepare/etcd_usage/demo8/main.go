package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

const KEY = `/cron/jobs/key1`
const VALUE = `value1`

// 使用OP操作取代kv.Put()、kv.Get()、kv.Delete()操作
func main() {
    var (
        config     clientv3.Config
        client     *clientv3.Client
        err        error
        kv         clientv3.KV
        opPut      clientv3.Op
        opResponse clientv3.OpResponse

        opGet clientv3.Op
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

    // 创建OP
    opPut = clientv3.OpPut(KEY, VALUE)
    // 执行OP
    if opResponse, err = kv.Do(context.TODO(), opPut); err != nil {
        fmt.Printf("Failed to do Put OP, err: %s", err.Error())
        return
    }

    // 打印写入的Revision
    fmt.Printf("写入的Revision: %d\n", opResponse.Put().Header.Revision)

    // 创建OP
    opGet = clientv3.OpGet(KEY)
    // 执行OP
    if opResponse, err = kv.Do(context.TODO(), opGet); err != nil {
        fmt.Printf("Failed to do Get OP, err: %s", err.Error())
        return
    }
    fmt.Printf("获取数据Revision: %d\n", opResponse.Get().Kvs[0].ModRevision)
    fmt.Printf("数据值：%s", opResponse.Get().Kvs[0].Value)
}
