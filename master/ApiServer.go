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
	// 单例对象
	G_apiServer *ApiServer
)

// 保存任务接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {

	// 任务保存到 etcd 中
}

func InitApiServer() error {
	// 配置路由
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 创建监听器
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort))

	if err != nil {
		return err
	}

	// 创建一个HTTP服务
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 启动服务
	go httpServer.Serve(listener)

	return nil
}
