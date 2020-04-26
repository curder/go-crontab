package common

import "encoding/json"

type Job struct {
    Name     string `json:"name"`     // 任务名
    Command  string `json:"command"`  // shell命令
    CronExpr string `json:"cron_expr"` // cron表达式
}

// 定义HTTP响应接口结构体
type Response struct {
    ErrorNumber int         `json:"error_number"`
    Message     string      `json:"message"`
    Data        interface{} `json:"data"`
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
