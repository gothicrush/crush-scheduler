package common

import "encoding/json"

// 定时任务
type Job struct {
	Name     string `json:"name"`     // 任务名
	Command  string `json:"command"`  //shell 命令
	CronExpr string `json:"cronExpr"` // cron 表达式
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
