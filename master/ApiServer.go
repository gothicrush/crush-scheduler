package master

import (
	"encoding/json"
	"github.com/gothicrush/crush-scheduler/common"
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
	// POST job={"name":"job1", "command":"echo hello", "cronExpr":"* * * * *"}

	// 任务保存到 etcd 中

	// 解析 POST 表单
	err := r.ParseForm()

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	// 获取表单中job字段
	postJob := r.PostForm.Get("myjob")

	// 反序列化job
	var job common.Job

	err = json.Unmarshal([]byte(postJob), &job)

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	// 保存到 etcd
	oldJob, err := G_jobManager.SaveJob(&job)

	// 返回正常应答
	resp, err := common.BuildResponse(0, "save success", oldJob)

	if err == nil {
		w.Write(resp)
		return
	}
}

// 删除任务接口
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	// Post /job/delete name=job1

	// 解析表单
	err := r.ParseForm()

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	// 获取删除任务名
	name := r.PostForm.Get("name")

	// 删除任务
	oldJob, err := G_jobManager.DeleteJob(name)

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	resp, _ := common.BuildResponse(0, "delete success", oldJob)
	w.Write(resp)
}

// 列出所有任务
func handleJobList(w http.ResponseWriter, r *http.Request) {

	var jobList []*common.Job

	jobList, err := G_jobManager.ListJob()

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	resp, _ := common.BuildResponse(0, "list success", jobList)
	w.Write(resp)
}

// 强制删除某个路由
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	//解析表单
	err := r.ParseForm()

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	// 获取杀死任务的名字
	name := r.PostForm.Get("name")

	// 杀死任务
	err = G_jobManager.KillJob(name)

	if err != nil {
		resp, _ := common.BuildResponse(-1, err.Error(), nil)
		w.Write(resp)
		return
	}

	// 正常响应
	resp, _ := common.BuildResponse(0, "kill success", nil)
	w.Write(resp)
}

func InitApiServer() error {

	// 配置路由
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	// 配置静态文件
	// 静态文件目录
	staticDir := http.Dir(G_config.WebRoot)
	staticHandler := http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

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
