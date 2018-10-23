package main

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {

	// 连接的配置
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立客户端连接
	client, err := clientv3.New(config)

	if err != nil {
		log.Println("clientv3.New error : ", err)
		return
	}

	// 建立KV客户端
	kvClient := clientv3.NewKV(client)

	// 使用PUT操作

	// 不带附加操作的PUT操作
	var putResp *clientv3.PutResponse

	putResp, err = kvClient.Put(context.TODO(), "/cron/jobs/job2", "world")

	if err != nil {
		log.Println("kvClient.Put error : ", err)
	}

	log.Println("ret:", putResp)
}
