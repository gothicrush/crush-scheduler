package main

import (
	"flag"
	"fmt"
	"github.com/gothicrush/crush-scheduler/master"
	"runtime"
	"time"
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

	// 启动任务管理器，与 etcd 服务器建立连接
	if err := master.InitJobManager(); err != nil {
		fmt.Println(err)
		return
	}

	// 启动API HTTP服务
	if err := master.InitApiServer(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
