package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

func main() {
    var (
        config      clientv3.Config
        client      *clientv3.Client
        err         error
        kv          clientv3.KV
        putResponse *clientv3.PutResponse
    )

    // 配置etcd
    config = clientv3.Config{
        Endpoints:   []string{"192.168.205.10:2379"}, // etcd集群配置
        DialTimeout: 5 * time.Second,
    }

    // 建立客户端连接
    if client, err = clientv3.New(config); err != nil {
        fmt.Printf("Failed to connect etcd server, err: %s", err.Error())
        return
    }
    defer client.Close()

    // 读写etcd的键值对
    kv = clientv3.NewKV(client)

    // 写入
    if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/key1", "value", clientv3.WithPrevKV()); err != nil {
        fmt.Printf("Failed to put key to etcd, err: %s", err.Error())
        return
    }

    fmt.Printf("Revision: %d \n", putResponse.Header.Revision) // 存储的版本
    if putResponse.PrevKv != nil {                             // 获取上一次存储的值
        fmt.Printf("PrevValue: %s \n", putResponse.PrevKv.Value)
    }
}
