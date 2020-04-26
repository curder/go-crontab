package main

import (
	"flag"
	"fmt"
	"github.com/curder/go-crontab/crontab/master"
	"runtime"
)

var (
	configFile string // 配置文件所在路径
)

func main() {
	var (
		err error
	)

	// 初始化命令行参数
	initArgs()

	// 加载配置
	if err = master.InitConfig(configFile); err != nil {
		goto ERR
	}
	// 初始化线程
	initEnv()

	// 启动http服务器
	if err = master.InitAPiServer(); err != nil {
		goto ERR
	}

	// 正常退出
	return

ERR:
	fmt.Printf("err: %s", err)
}

// 解析命令行参数
func initArgs() {
	// master -config ./master.json
	flag.StringVar(&configFile, "config", "./master.json", "配置master.json文件路径")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
