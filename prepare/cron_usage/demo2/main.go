package main

import (
    "fmt"
    "github.com/gorhill/cronexpr"
    "time"
)

func main() {
    // 需要有调度协程，它定时检查所有的Cron任务，有任务过期了就执行
    var (
        cronJob       *CronJob
        expr          *cronexpr.Expression
        now           time.Time
        scheduleTable map[string]*CronJob // key: 任务的名称
    )
    // 初始化调度表
    scheduleTable = make(map[string]*CronJob)

    // 当前时间
    now = time.Now()

    // 定义定时任务
    expr = cronexpr.MustParse("*/5 * * * * * *") // 每隔5s执行
    cronJob = &CronJob{
        expr:    expr,
        nexTime: expr.Next(now),
    }

    // 任务注册到调度表，命名为job1
    scheduleTable["job1"] = cronJob

    // 定义定时任务
    expr = cronexpr.MustParse("*/10 * * * * * *") // 每隔10s执行
    cronJob = &CronJob{
        expr:    expr,
        nexTime: expr.Next(now),
    }

    // 任务注册到调度表，命名为job2
    scheduleTable["job2"] = cronJob

    // 启动调度携程
    go func() {
        var (
            jobName string
            cronJob *CronJob
        )
        // 定时检查一下任务调度表
        for {
            now = time.Now()

            for jobName, cronJob = range scheduleTable {
                // 判断任务是否过期
                if cronJob.nexTime.Before(now) || cronJob.nexTime.Equal(now) {
                    // 启动协程,执行任务
                    go func(jobName string) {
                        fmt.Println("当前执行的任务是: ", jobName)
                    }(jobName)

                    // 计算下一次调度时间
                    cronJob.nexTime = cronJob.expr.Next(now)
                    fmt.Println(jobName, " 下次调度的时间是： ", cronJob.nexTime)
                }
            }
            // 休眠 100 ms
            select {
            case <-time.NewTimer(100 * time.Millisecond).C:
            }
        }
    }()

    // 为了避免主协程退出，休眠100s
    time.Sleep(100 * time.Second)
}

type CronJob struct {
    expr    *cronexpr.Expression
    nexTime time.Time
}
