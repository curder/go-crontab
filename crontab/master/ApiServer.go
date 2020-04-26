package master

import (
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

	defer httpServer.Close()

	// 赋值单例
	GApiServer = &ApiServer{httpServer: httpServer}

	// 启动服务
	go httpServer.Serve(listener)

	return
}

func handleJobSave(w http.ResponseWriter, r *http.Request) {
	//
}
