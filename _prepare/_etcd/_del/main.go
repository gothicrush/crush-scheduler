package main

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {

	// 连接配置
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	client, err := clientv3.New(config)

	if err != nil {
		log.Println("clientv3.New error : ", err)
		return
	}

	// 建立KV客户端
	kvClient := clientv3.NewKV(client)

	// 使用DEL操作

	// 不带任何额外参数
	delResp, err := kvClient.Delete(context.TODO(), "/cron/jobs/job1")

	if err != nil {
		log.Println("kvClient.Delete 1 error : ", err)
		return
	}

	log.Println("delete 1 : ", delResp.Header)

	// 带 --prefix 额外参数
	delResp, err = kvClient.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())

	if err != nil {
		log.Println("kvClient.Delete 2 error : ", err)
		return
	}

	log.Println("delete 2 : ", delResp.Header)

	// 查看是否删除成功

	getResp, err := kvClient.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())

	if err != nil {
		log.Println("kvClient.Get error : ", err)
	}

	log.Println("Count : ", getResp.Count)
}
