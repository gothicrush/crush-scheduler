package main

import (
	"context"
	"fmt"
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

	// 建立Lease客户端
	leaseClient := clientv3.NewLease(client)

	// 创建一个10秒的租约，并获取其ID
	leaseOne, err := leaseClient.Grant(context.TODO(), 10)

	if err != nil {
		log.Println("leaseClient.Grant 1 error : ", err)
		return
	}

	leaseOneID := leaseOne.ID

	// 伴随leaseOne放入一个值
	kvClient.Put(context.TODO(), "/cron/lock/job1", "hello", clientv3.WithLease(leaseOneID))

	// 每2秒检查一次， /cron/lock/job1 是否过期
	var getResp *clientv3.GetResponse
	for {
		getResp, err = kvClient.Get(context.TODO(), "/cron/lock/job1", clientv3.WithCountOnly())

		if err != nil {
			log.Println("kvClient.Get error : ", err)
			return
		}

		log.Println("Count : ", getResp.Count)
		if getResp.Count == 0 {
			break
		}

		select {
		case <-time.NewTimer(2 * time.Second).C:
		}
	}

	fmt.Println("------------------------------------------------")

	// 创建一个新的租约，5秒，获取其 ID

	leaseTwo, err := leaseClient.Grant(context.TODO(), 5)

	if err != nil {
		log.Println("leaseClient.Grant 2 error : err")
		return
	}

	leaseTwoID := leaseTwo.ID

	// 开启自动续约
	leaseClient.KeepAlive(context.TODO(), leaseTwoID)

	// 带着租约放入一个值
	kvClient.Put(context.TODO(), "/cron/lock/job2", "world", clientv3.WithLease(leaseTwoID))

	// 每2秒读取一次 /cron/lock/job2 ，观察10秒后还存不存在
	times := 0
	for times < 10 {
		getResp, err = kvClient.Get(context.TODO(), "/cron/lock/job2")

		if err != nil {
			log.Println("kvClient.Get 2 error : ", err)
			return
		}

		log.Println(getResp.Count)

		select {
		case <-time.NewTimer(2 * time.Second).C:
		}

		times++
	}

}
