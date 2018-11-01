package main

import (
	"flag"
	"fmt"
	"github.com/gothicrush/crush-scheduler/master"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "传入 master.json")
	flag.Parse()
}

func main() {

	// 设置线程数目
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 处理命令行参数
	initArgs()

	// 加载配置
	if err := master.InitConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化服务发现
	if err := master.InitWorkerManager(); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化日志管理器
	if err := master.InitLogManager(); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化任务管理器
	if err := master.InitJobManager(); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化API HTTP服务
	if err := master.InitApiServer(); err != nil {
		fmt.Println(err)
		return
	}

	var temp chan int = make(chan int)
	for {
		select {
		case <-temp:
		}
	}
}
