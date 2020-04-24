package main

import (
    "fmt"
    "github.com/gorhill/cronexpr"
    "time"
)

func main() {
    var (
        expr     *cronexpr.Expression
        err      error
        now      time.Time
        nextTime time.Time
    )

    // 每几秒 那一年(2018-2099) 哪一分钟(0-59) 哪小时(0-23) 哪天(1-31) 哪月(1-12) 星期几(0-6)

    // 每隔5分钟执行一次
    if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
        fmt.Println(err)
        return
    }

    // 当前时间
    now = time.Now()

    // 下次调用时间
    nextTime = expr.Next(now)

    // 等待定时器超时
    time.AfterFunc(nextTime.Sub(now), func() {
        fmt.Println("被调度:" , nextTime)
    })

    time.Sleep(5 * time.Second)
}
