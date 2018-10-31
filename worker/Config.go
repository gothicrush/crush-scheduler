package worker

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	EtcdEndPoints         []string `json:"etcdEndPoints"`
	EtcdTimeout           int      `json:"etcdTimeout"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
	JobLogBatchSize       int      `json:"jobLogBatchSize"`
	JobLogCommitTimeout   int      `json:"jobLogCommitTimeout"`
}

// 单例
var (
	G_config *Config
)

func InitConfig(filename string) error {

	// 读取配置文件
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	// json反序列化
	var conf Config

	err = json.Unmarshal(content, &conf)

	if err != nil {
		return err
	}

	// 单例赋值
	G_config = &conf

	return nil
}
