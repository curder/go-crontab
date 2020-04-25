package main

import (
    "context"
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "time"
)

const KEY = `/cron/jobs/key10`
const VALUE = `value10`

// lease锁自动过期、OP操作和txn事务：if else then
func main() {
    var (
        config                 clientv3.Config
        client                 *clientv3.Client
        err                    error
        lease                  clientv3.Lease
        leaseGrantResponse     *clientv3.LeaseGrantResponse
        leaseID                clientv3.LeaseID
        leaseKeepAliveResponse <-chan *clientv3.LeaseKeepAliveResponse
        keepResponse           *clientv3.LeaseKeepAliveResponse
        ctx                    context.Context
        cancelFunc             context.CancelFunc
        kv                     clientv3.KV
        txn                    clientv3.Txn
        txnResponse            *clientv3.TxnResponse
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

    // 上锁（创建租约并自动续租，拿着租约去抢占一个key）
    // 创建租约
    lease = clientv3.NewLease(client)
    if leaseGrantResponse, err = lease.Grant(context.TODO(), 10); err != nil {
        fmt.Printf("Failed to grant response, err: %s", err.Error())
        return
    }
    leaseID = leaseGrantResponse.ID

    // 准备一个用于自动续租的context
    ctx, cancelFunc = context.WithCancel(context.TODO())

    // 确保函数退出后，自动续租停止
    defer cancelFunc()
    defer lease.Revoke(context.TODO(), leaseID)

    // 10 S 后自动续租
    if leaseKeepAliveResponse, err = lease.KeepAlive(ctx, leaseID); err != nil {
        fmt.Printf("Failed to keepalive, err: %s", err.Error())
        return
    }
    go func() {
    END:
        for {
            select {
            case keepResponse = <-leaseKeepAliveResponse:
                if keepResponse == nil {
                    fmt.Println("租约失效了")
                    goto END
                } else {
                    fmt.Println("收到自动续约应答：", keepResponse.ID)
                }
            }
        }
    }()

    // 如果不存在Key，设置，否则抢key失败
    kv = clientv3.NewKV(client)
    // 创建事务
    txn = kv.Txn(context.TODO())

    txn.If(clientv3.Compare(clientv3.CreateRevision(KEY), "=", 0)). // 如果不存在key
        Then(clientv3.OpPut(KEY, VALUE, clientv3.WithLease(leaseID))). // 设置key
        Else(clientv3.OpGet(KEY)) // 抢锁失败

    if txnResponse, err = txn.Commit(); err != nil { // 提交锁
        fmt.Printf("Failed to commit txt, err: %s\n", err.Error())
        return
    }

    if !txnResponse.Succeeded {
        fmt.Printf("锁被占用：%s\n", txnResponse.Responses[0].GetResponseRange().Kvs[0].Value)
        return
    }

    // 处理业务

    fmt.Println("处理任务")
    time.Sleep(time.Second * 5)

    // 释放锁（取消自动续租，释放租约「立即删除KV」）
    // defer 会自动取消租约，关联的KV就被删除了

}
