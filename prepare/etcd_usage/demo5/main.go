package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "go.etcd.io/etcd/mvcc/mvccpb"
    "time"
)

func main() {
    var (
        config clientv3.Config
        client *clientv3.Client
        err    error
        kv     clientv3.KV

        deleteResponse *clientv3.DeleteResponse
        kvPair         *mvccpb.KeyValue
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

    // 删除KV
    if deleteResponse, err = kv.Delete(context.TODO(), "/cron/jobs/key1", clientv3.WithPrevKV()); err != nil { // 删除多个key： kv.Delete(context.TODO(), "/cron/jobs", clientv3.WithPrefix())
        fmt.Printf("Failed to delete key, err: %s", err.Error())
        return
    }

    // 如果删除成功，查看被删除前的kv
    if len(deleteResponse.PrevKvs) != 0 {
        for _, kvPair = range deleteResponse.PrevKvs {
            fmt.Printf("删除了：%s, %s", kvPair.Key, kvPair.Value)
        }
    } else {
        fmt.Println("没有任何key需要被删除")
    }
}
