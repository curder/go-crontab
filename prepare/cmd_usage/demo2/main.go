package main

import (
    "context"
    "fmt"
    "os/exec"
    "time"
)

// 杀死进程
func main() {
    // 在协程中执行一个cmd，并执行2秒：sleep 2; echo hello; 并在1s的时候杀死cmd

    var (
        ctx        context.Context
        cancelFunc context.CancelFunc
        resultChan chan *result
        res        *result
    )

    // 创建一个输出结果队列
    resultChan = make(chan *result, 1000)

    // 创建一个上下文
    ctx, cancelFunc = context.WithCancel(context.TODO())

    go func() {
        var (
            cmd    *exec.Cmd
            output []byte
            err    error
        )
        cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2; echo hello;")
        if output, err = cmd.CombinedOutput(); err != nil {
            fmt.Printf("cmd.CombinedOutput err: %s", err.Error())
            return
        }

        // 构造输出结果
        resultChan <- &result{
            output: output,
            err:    err,
        }

    }()

    time.Sleep(1 * time.Second)
    cancelFunc() // 取消子进程

    // 在main协程等待子协程退出后打印子协程输出内容
    res = <-resultChan
    fmt.Printf("Output: %s, Err: %v", res.output, res.err)
}

type result struct {
    output []byte
    err    error
}
