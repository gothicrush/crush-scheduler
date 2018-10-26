package main

import (
	"flag"
	"fmt"
	"github.com/gothicrush/crush-scheduler/master"
	"runtime"
)

// 初始化线程数目
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./main.json", "传入 main.json")
	flag.Parse()
}

func main() {

	// 设置线程数目
	initEnv()

	// 处理命令行参数
	initArgs()

	// 加载配置
	if err := master.InitConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}

	// 启动任务管理器
	if err := master.InitJobMgr(); err != nil {
		fmt.Println(err)
		return
	}

	// 启动API HTTP服务
	if err := master.InitApiServer(); err != nil {
		fmt.Println(err)
		return
	}
}
