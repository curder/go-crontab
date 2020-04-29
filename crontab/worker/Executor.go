package worker

import (
    "context"
    "github.com/curder/go-crontab/crontab/common"
    "os/exec"
    "time"
)

type Executor struct {
}

var (
    GExecutor *Executor
)

func (e *Executor) ExecuteJob(info *common.JobExecuteInfo) {
    go func() {
        var (
            cmd    *exec.Cmd
            output []byte
            err    error
            result *common.JobExecuteResult
        )

        // 任务执行结果
        result = &common.JobExecuteResult{
            ExecuteInfo: info,
            Output:      make([]byte, 0),
        }

        result.StartTime = time.Now()

        // 执行shell命令
        cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.CronExpr)

        // 执行并捕获输出
        output, err = cmd.CombinedOutput()

        // 记录任务结束时间
        result.EndTime = time.Now()
        result.Output = output
        result.Err = err

        // 任务执行完成后，把执行的结构返回给Scheduler
        GScheduler.PushJobResult(result)
    }()
}

func InitExecutor() (err error) {
    GExecutor = &Executor{

    }

    return
}
