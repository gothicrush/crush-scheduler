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

	defer client.Close()

	// 创建 op 操作
	// OpPut
	opPut := clientv3.OpPut("/op/test", "/op/test/put")
	// OpGet
	opGet := clientv3.OpGet("/op/test")
	// OpDelete
	opDelete := clientv3.OpDelete("/op/test")

	// 创建kv客户端
	kvClient := clientv3.NewKV(client)

	// 执行opPut操作
	opPutResponse, err := kvClient.Do(context.TODO(), opPut)

	if err != nil {
		log.Println("kvClient.Do.opPut error : ", err)
		return
	}

	log.Println(opPutResponse.Put().Header.Revision)

	// 执行opGet操作
	opGetResponse, err := kvClient.Do(context.TODO(), opGet)

	if err != nil {
		log.Println("kvClient.Do.opGet error : ", err)
		return
	}

	log.Println(opGetResponse.Get().Count)

	// 执行opDelete操作
	opDeleteResponse, err := kvClient.Do(context.TODO(), opDelete)

	if err != nil {
		log.Println("kvClient.Do.opDelete error : ", err)
		return
	}

	log.Println(opDeleteResponse.Del().Header.Revision)

	// 执行opGet操作
	opGetResponse, err = kvClient.Do(context.TODO(), opGet)

	if err != nil {
		log.Println("kvClient.Do.opGet error 2 : ", err)
		return
	}

	log.Println(opGetResponse.Get().Count)
}
