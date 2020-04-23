package main

import (
    "fmt"
    "os/exec"
)

// 使用exec执行系统命令

func main() {
    var (
        cmd    *exec.Cmd
        output []byte
        err    error
    )

    // 生成执行cmd
    cmd = exec.Command("/bin/bash", "-c", "sleep 2; ls -l")

    // 执行命令，捕获子进程输出
    if output, err = cmd.CombinedOutput(); err != nil {
        fmt.Println(err)
        return
    }

    // 打印子进程输出
    fmt.Printf("%s", output)
}
