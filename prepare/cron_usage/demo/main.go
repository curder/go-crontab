package main

import (
    "fmt"
    "github.com/gorhill/cronexpr"
)

func main() {
    var (
        expr *cronexpr.Expression
        err   error
    )

    // 哪一分钟(0-59) 哪小时(0-23) 哪天(1-31) 哪月(1-12) 星期几(0-6)

    // 每分钟执行一次
    if expr, err = cronexpr.Parse("* * * * *"); err != nil {
        fmt.Println(err)
        return
    }

    // 每隔5分钟执行一次
    if expr, err = cronexpr.Parse("*/5 * * * *"); err != nil {
        fmt.Println(err)
        return
    }
    expr = expr
}
