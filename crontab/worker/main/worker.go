package main

import (
    "flag"
    "fmt"
    "github.com/curder/go-crontab/crontab/worker"
    "runtime"
    "time"
)

var (
    configFile string // 配置文件所在路径
)

func main() {
    // 解析命令行参数
    var (
        err error
    )

    // 初始化命令行参数
    initArgs()

    // 初始化线程
    initEnv()

    // 加载配置
    if err = worker.InitConfig(configFile); err != nil {
        goto ERR
    }

    // 启动任务调度器
    if err = worker.InitScheduler(); err != nil {
        goto ERR
    }

    // 任务etcd管理器
    if err = worker.InitJobMgr(); err != nil {
        goto ERR
    }

    // 正常退出
    for {
        time.Sleep(1 * time.Second)
    }

ERR:
    fmt.Printf("err: %s", err)
}

// 解析命令行参数
func initArgs() {
    // worker -config ./worker.json
    flag.StringVar(&configFile, "config", "./worker.json", "配置worker.json文件路径")
    flag.Parse()
}

func initEnv() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}
