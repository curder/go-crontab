package common

import (
    "encoding/json"
    "strings"
)

type Job struct {
    Name     string `json:"name"`      // 任务名
    Command  string `json:"command"`   // shell命令
    CronExpr string `json:"cron_expr"` // cron表达式
}

// 定义HTTP响应接口结构体
type Response struct {
    ErrorNumber int         `json:"error_number"`
    Message     string      `json:"message"`
    Data        interface{} `json:"data"`
}

type JobEvent struct {
    EventType int // SAVE、Delete
    Job       *Job
}

// 响应方法
func BuildResponse(errorNumber int, message string, data interface{}) (response []byte, err error) {
    // 定义一个response
    var (
        resp Response
    )

    resp.ErrorNumber = errorNumber
    resp.Message = message
    resp.Data = data

    // 反序列化json
    response, err = json.Marshal(resp)
    return
}

// 反序列化Job
func UnpackJob(value []byte) (ret *Job, err error) {
    var (
        job *Job
    )

    job = &Job{}
    if err = json.Unmarshal(value, job); err != nil {
        return
    }

    ret = job

    return
}

// 从 etcd 的key中提取任务名， /cron/jobs/jobName => jobName
func ExtraJobName(jobKey string) string {
    return strings.TrimPrefix(jobKey, JobSaveDir)
}

// 任务变化事件有2种，1.更新 2. 删除
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
    return &JobEvent{
        EventType: eventType,
        Job:       job,
    }
}
