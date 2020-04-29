package worker

import (
    "fmt"
    "github.com/curder/go-crontab/crontab/common"
    "time"
)

type Scheduler struct {
    jobEventChan      chan *common.JobEvent              // etcd任务事件队列
    jobPlanTable      map[string]*common.JobSchedulePlan // 任务调度计划池
    jobExecutingTable map[string]*common.JobExecuteInfo  // 任务执行池
    jobResultChan     chan *common.JobExecuteResult      // 任务结果队列
}

var (
    GScheduler *Scheduler
)

// 处理事件
func (s *Scheduler) handlerJobEvent(jobEvent *common.JobEvent) {
    var (
        jobExists       bool
        jobSchedulePlan *common.JobSchedulePlan
        err             error
    )
    switch jobEvent.EventType {
    case common.JobEventSave: // 保存任务事件
        if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
            return
        }
        s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
    case common.JobEventDelete: // 删除任务事件
        if jobSchedulePlan, jobExists = s.jobPlanTable[jobEvent.Job.Name]; jobExists { // 如果任务存在
            delete(s.jobPlanTable, jobEvent.Job.Name)
        }
    }
}

// 尝试执行任务
func (s *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
    // 执行的任务可能运行很久，1分钟会被调度60次，但是只能执行一次，防止并发

    var (
        jobExecuteInfo *common.JobExecuteInfo
        jobExecuting   bool
    )
    // 如果任务正在执行，跳过本次调度
    if jobExecuteInfo, jobExecuting = s.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
        fmt.Println("任务尚未退出执行，跳过此次执行", jobPlan.Job.Name)
        return
    }

    // 构建执行状态信息
    jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

    // 保存执行状态
    s.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

    // TODO 执行任务
    GExecutor.ExecuteJob(jobExecuteInfo)
    // fmt.Println("执行任务：", jobPlan.Job.Name, "计划执行时间：", jobExecuteInfo.PlanTime, "真正执行时间：", jobExecuteInfo.RealTime)
}

// 重新计算任务调度状态
func (s *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
    var (
        jobPlan  *common.JobSchedulePlan
        now      = time.Now()
        nearTime *time.Time
    )

    if len(s.jobPlanTable) == 0 { // 如果任务表中不存在任务
        time.Sleep(1 * time.Second) // 休眠1秒
        return
    }

    // 遍历所有任务，过期的任务立即执行
    for _, jobPlan = range s.jobPlanTable {
        if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
            // 尝试执行任务
            s.TryStartJob(jobPlan)

            fmt.Println("执行任务", jobPlan.Job.Name)
            jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
        }

        // 统计最近的要过期的任务时间
        if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
            nearTime = &jobPlan.NextTime
        }

        // 下次调度间隔
        scheduleAfter = (*nearTime).Sub(now)

    }
    return
}

// 调度协程
func (s *Scheduler) loop() {
    var (
        jobEvent      *common.JobEvent
        scheduleAfter time.Duration
        scheduleTimer *time.Timer
        jobResult     *common.JobExecuteResult
    )

    // 初始化任务调度
    scheduleAfter = s.TrySchedule()

    // 延时调度器
    scheduleTimer = time.NewTimer(scheduleAfter)

    // 定时任务
    for {
        select {
        case jobEvent = <-s.jobEventChan: // 监听任务变化事件
            s.handlerJobEvent(jobEvent)
        case <-scheduleTimer.C: // 最近的任务到期
        case jobResult = <-s.jobResultChan: // 任务执行结果
            s.handlerJobResult(jobResult)
        }
        s.TrySchedule()                    // 调取任务
        scheduleTimer.Reset(scheduleAfter) // 重置调度间隔
    }
}

// 推送任务变化事件

func (s *Scheduler) PushJobEvent(event *common.JobEvent) {
    s.jobEventChan <- event
}

// 初始调度器
func InitScheduler() (err error) {
    GScheduler = &Scheduler{
        jobEventChan:      make(chan *common.JobEvent, 1000),
        jobPlanTable:      make(map[string]*common.JobSchedulePlan),
        jobExecutingTable: make(map[string]*common.JobExecuteInfo),
        jobResultChan:     make(chan *common.JobExecuteResult, 1000),
    }

    // 启动调度协程
    go GScheduler.loop()

    return
}

// 回传任务执行结果
func (s *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
    s.jobResultChan <- jobResult
}

// 处理任务结果
func (s *Scheduler) handlerJobResult(result *common.JobExecuteResult) {
    delete(s.jobExecutingTable, result.ExecuteInfo.Job.Name)

    fmt.Println("任务处理完成", result.ExecuteInfo.Job.Name, string(result.Output))
}
