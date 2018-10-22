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

	// 不带任何额外参数
	getResp, err := kvClient.Get(context.TODO(), "/cron/jobs/job2")

	if err != nil {
		log.Println("kvClient.Get error : ", err)
		return
	}

	log.Println("不带参数 ret : Count: ", getResp.Count, " ret: ", getResp.Kvs)

	// 带 prefix 参数
	getResp, err = kvClient.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())

	if err != nil {
		log.Println("kvClient.Get error : ", err)
		return
	}

	log.Println("带WithPrefix参数 ret : Count: ", getResp.Count, " ret: ", getResp.Kvs)

}
