package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name     string `json:"name"`     // 任务名
	Command  string `json:"command"`  //shell 命令
	CronExpr string `json:"cronExpr"` // cron 表达式
}

// 反序列化job
func UnpackJob(value []byte) (*Job, error) {

	job := Job{}

	err := json.Unmarshal(value, &job)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

// 任务变化事件有：1.更新事件，2.删除事件
type JobEvent struct {
	EventType int //SAVE,DELETE
	Job       *Job
}

// 构建任务事件
func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// 任务执行状态
type JobExecuteInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 理论上的调度时间
	ReadTime   time.Time          // 实际调度时间
	CancelCtx  context.Context    // 用于取消任务的context
	CancelFunc context.CancelFunc // 用于取消命令执行的cancel函数
}

// 构造任务执行信息
func BuildJobExecuteInfo(jobPlan *JobSchedulePlan) *JobExecuteInfo {
	jobExecuteInfo := &JobExecuteInfo{
		Job:      jobPlan.Job,
		PlanTime: jobPlan.NextTime,
		ReadTime: time.Now(),
	}

	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())

	return jobExecuteInfo
}

// 任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 // 要调度的任务
	Expr     *cronexpr.Expression // cron表达式
	NextTime time.Time            // 下次调度时间
}

// 构造任务执行计划
func BuildJobSchedulePlan(job *Job) (*JobSchedulePlan, error) {
	// 解析cron表达式
	expr, err := cronexpr.Parse(job.CronExpr)

	if err != nil {
		return nil, err
	}

	// 生成调度计划表
	jobSchedulePlan := &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return jobSchedulePlan, nil
}

// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          //命令输出
	Err         error           // 命令执行错误原因
	StartTime   time.Time       //启动时间
	EndTime     time.Time       // 结束时间
}

////////////////////////////////////////

// HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 应答方法
func BuildResponse(errno int, msg string, data interface{}) ([]byte, error) {
	var response Response

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	ret, err := json.Marshal(response)

	if err != nil {
		return nil, err
	}

	return ret, nil
}

///////////////////////////////////////////////////

// 任务执行日志
type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`           //任务名字
	Command      string `json:"command" bson:"command"`           //脚本命令
	Err          string `json:"err" bson:"err"`                   //错误原因
	Output       string `json:"output" bson:"output"`             //脚本输出
	PlanTime     int64  `json:"planTime" bson:"planTime"`         //计划开始时间
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"` //实际调度时间
	StartTime    int64  `json:"startTime" bson:"startTime"`       //任务开始时间
	EndTime      int64  `json:"endTime" bson:"endTime"`           //任务结束时间
}

// 任务日志批次
type LogBatch struct {
	Logs []interface{} // 多条日志
}

// 任务日志过滤条件
type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

// 任务日志排序规则
type SortLogByStartTime struct {
	StartOrder int `bson:"startTime"`
}

//////////////////////////////////////

// 从etcd的job中提取任务名
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// 从etcd的killer中提取任务名
func ExtractKillerName(killerKey string) string {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

// 从etcd的worker中提取worker的IP
func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}
