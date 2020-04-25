package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

func main() {
    var (
        config clientv3.Config
        client *clientv3.Client
        err    error
        kv     clientv3.KV

        getResponse *clientv3.GetResponse
    )
    config = clientv3.Config{
        Endpoints:   []string{"192.168.205.10:2379"},
        DialTimeout: 5 * time.Second,
    }

    if client, err = clientv3.New(config); err != nil {
        fmt.Printf("Failed to connect etcd err: %s", err.Error())
        return
    }

    defer client.Close()

    // 用于读写etcd的键值对
    kv = clientv3.NewKV(client)

    // 写入两个不同的key，分别是key1和key2
    if _, err = kv.Put(context.TODO(), "/cron/jobs/key1", "value1"); err != nil {
        fmt.Printf("Failed to put key to etcd server, err; %s", err.Error())
        return
    }
    if _, err = kv.Put(context.TODO(), "/cron/jobs/key2", "value2"); err != nil {
        fmt.Printf("Failed to put key to etcd server, err; %s", err.Error())
        return
    }

    // 获取 /cron/jobs/ 为前缀的所有key
    if getResponse, err = kv.Get(context.TODO(), "/cron/jobs/key", clientv3.WithPrefix()); err != nil {
        fmt.Printf("Failed to get response err: %s", err.Error())
        return
    }

    // 打印所有kvs
    fmt.Println(getResponse.Kvs)

}
