package worker

import "github.com/curder/go-crontab/crontab/common"

type Scheduler struct {
    jobEventChan chan *common.JobEvent              // etcd任务事件队列
    jobPlanTable map[string]*common.JobSchedulePlan // 任务调度计划池
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

// 调度协程
func (s *Scheduler) loop() {
    var (
        jobEvent *common.JobEvent
    )
    // 定时任务
    for {
        select {
        case jobEvent = <-s.jobEventChan:
            s.handlerJobEvent(jobEvent)
        }
    }
}

// 推送任务变化事件

func (s *Scheduler) PushJobEvent(event *common.JobEvent) {
    s.jobEventChan <- event
}

// 初始调度器
func initScheduler() {
    GScheduler = &Scheduler{
        jobEventChan: make(chan *common.JobEvent, 1000),
        jobPlanTable: make(map[string]*common.JobSchedulePlan),
    }

    // 启动调度协程
    go GScheduler.loop()

    return
}
