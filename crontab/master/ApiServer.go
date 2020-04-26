package master

import (
    "encoding/json"
    "fmt"
    "github.com/curder/go-crontab/crontab/common"
    "net"
    "net/http"
    "strconv"
    "time"
)

type ApiServer struct {
    httpServer *http.Server
}

var (
    GApiServer *ApiServer // 单例对象
)

// 初始化服务
func InitAPiServer() (err error) {
    var (
        mux        *http.ServeMux
        listener   net.Listener
        httpServer *http.Server
    )

    // 配置路由
    mux = http.NewServeMux()
    mux.HandleFunc("/job/save", handleJobSave)

    // 启动TCD监听
    if listener, err = net.Listen("tcp", ":"+strconv.Itoa(GConfig.APiPort)); err != nil {
        return
    }

    // 创建http服务器
    httpServer = &http.Server{
        Handler:      mux,
        ReadTimeout:  time.Duration(GConfig.ApiReadTimeout) * time.Millisecond,
        WriteTimeout: time.Duration(GConfig.ApiWriteTimeout) * time.Millisecond,
    }

    // defer httpServer.Close()

    // 赋值单例
    GApiServer = &ApiServer{httpServer: httpServer}

    // 启动服务
    go httpServer.Serve(listener)

    return
}

// 任务保存接口 POST job={"name": "jobName", "command": "echo hello", "cronExpr", "* * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
    var (
        err     error
        postJob string
        job     common.Job
        oldJob  *common.Job
        bytes   []byte
    )
    // 解析到POST表单
    if err = r.ParseForm(); err != nil {
        goto ERR
    }

    // 获取表单job字段
    postJob = r.PostForm.Get("job")

    // 反序列化Job
    if err = json.Unmarshal([]byte(postJob), &job); err != nil {
        goto ERR
    }

    // fmt.Printf("%#v", job)
    // 保存到etcd
    if oldJob, err = GJobMgr.SaveJob(&job); err != nil {
        goto ERR
    }

    // 返回正常响应
    if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
        _, _ = w.Write(bytes)
    }
fmt.Println(string(bytes))
    return

ERR:
    // 返回异常响应
    if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
        _, _ = w.Write(bytes)
    }
}
