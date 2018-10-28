package common

import (
	"encoding/json"
	"strings"
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

// 从etcd的key中提取任务名
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// 任务变化事件有：1.更新事件，2.删除事件
type JobEvent struct {
	EventType int //SAVE,DELETE
	job       *Job
}

func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		job:       job,
	}
}

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
