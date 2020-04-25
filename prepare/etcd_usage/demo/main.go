package main

import (
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

// 下载etcd 的go客户端命令：go get -v -u go.etcd.io/etcd/clientv3
func main() {
    var (
        config clientv3.Config
        client *clientv3.Client
        err    error
    )

    // 客户端连接配置
    config = clientv3.Config{
        Endpoints:   []string{"192.168.205.10:2379"},  // etcd服务地址
        DialTimeout: 5 * time.Second, // 建立连接的超时时间
    }

    // 建立连接
    if client, err = clientv3.New(config); err != nil {
        fmt.Printf("Failed to connect etcd: %s", err.Error())
        return
    }

    // 延迟关闭连接
    defer client.Close()

    client = client
}
