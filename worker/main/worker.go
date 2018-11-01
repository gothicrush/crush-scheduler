package main

import (
	"flag"
	"fmt"
	"github.com/gothicrush/crush-scheduler/worker"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "传入 worker.json")
	flag.Parse()
}

func main() {

	// 设置线程数目
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 处理命令行参数
	initArgs()

	// 加载配置
	if err := worker.InitConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}

	// 服务注册
	if err := worker.InitRegister(); err != nil {
		fmt.Println(err)
	}

	// 启动日志协程
	if err := worker.InitLogSink(); err != nil {
		fmt.Println(err)
		return
	}

	// 启动执行器
	if err := worker.InitExecutor(); err != nil {
		fmt.Println(err)
		return
	}

	// 启动调度器
	worker.InitScheduler()

	// 启动任务管理器，与 etcd 服务器建立连接
	if err := worker.InitJobManager(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
