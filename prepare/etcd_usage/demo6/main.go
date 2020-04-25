package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

// etcd的kv自动过期和续租
func main() {
    var (
        config             clientv3.Config
        client             *clientv3.Client
        err                error
        lease              clientv3.Lease
        leaseID            clientv3.LeaseID
        leaseGrantResponse *clientv3.LeaseGrantResponse
        kv                 clientv3.KV
        putResponse        *clientv3.PutResponse
        getResponse        *clientv3.GetResponse
        aliveResponseChan  <-chan *clientv3.LeaseKeepAliveResponse
        keepAliveResponse  *clientv3.LeaseKeepAliveResponse
    )

    // 客户端连接配置
    config = clientv3.Config{
        Endpoints:   []string{"192.168.205.10:2379"},
        DialTimeout: 5 * time.Second,
    }

    // 建立连接
    if client, err = clientv3.New(config); err != nil {
        fmt.Printf("Failed to connect etcd err: %s", err.Error())
        return
    }

    defer client.Close() // 延迟关闭客户端

    // 申请租约 lease
    lease = clientv3.NewLease(client)

    // 申请一个10s的租约
    if leaseGrantResponse, err = lease.Grant(context.TODO(), 10); err != nil {
        fmt.Printf("Failed to grant lease, err: %s", err.Error())
        return
    }

    defer lease.Close() // 延迟关闭租约

    // 获取租约ID
    leaseID = leaseGrantResponse.ID

    // 自动续租
    if aliveResponseChan, err = lease.KeepAlive(context.TODO(), leaseID); err != nil {
        fmt.Printf("Failed to keep alive lease, err: %s", err.Error())
        return
    }

    // 处理续约应答的协程
    go func() {
    END:
        for {
            select {
            case keepAliveResponse = <-aliveResponseChan:
                if keepAliveResponse == nil {
                    fmt.Printf("租约已失效")
                    goto END
                } else {
                    fmt.Println("收到自动续约应答：", keepAliveResponse.ID)
                }
            }
        }
    }()

    // 获得KV对象
    kv = clientv3.NewKV(client)

    // Put一个KV，让它与租约关联起来，从而实现10s后自动过期
    if putResponse, err = kv.Put(context.TODO(), "/cron/jobs/key1", "value1", clientv3.WithLease(leaseID)); err != nil {
        fmt.Printf("Failed to put a lease grant, err: %s", err.Error())
        return
    }

    fmt.Printf("KV写入成功，租约ID为：%d\n", putResponse.Header.Revision)

    for {
        if getResponse, err = client.Get(context.TODO(), "/cron/jobs/key1"); err != nil {
            fmt.Printf("Failed to get, err: %s\n", err.Error())
            return
        }

        if getResponse.Count == 0 {
            fmt.Println("目标KV过期了")
            break
        }
        fmt.Printf("KV还没过期：%s\n", getResponse.Kvs)
        time.Sleep(2 * time.Second)
    }
}
