package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
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

	// 启动协程每秒添加一个值又删除一个值
	go func() {
		for {
			kvClient.Put(context.TODO(), "/watch", "ee")
			kvClient.Delete(context.TODO(), "/watch")

			select {
			case <-time.NewTimer(1 * time.Second).C:
			}
		}
	}()

	// 获取 Revision
	getResponse, err := kvClient.Get(context.TODO(), "/watch")

	if err != nil {
		log.Println("kvClient.Get error : ", err)
	}

	revision := getResponse.Header.Revision
	revision++

	// 创建 watcher
	watcher := clientv3.NewWatcher(client)

	// 获取监听管道
	watchResponseChannel := watcher.Watch(context.TODO(), "/watch", clientv3.WithRev(revision))

	// 消费管道内容

	for {
		select {
		case watchResponse := <-watchResponseChannel:
			for _, ev := range watchResponse.Events {
				switch ev.Type {
				case mvccpb.PUT:
					fmt.Println("修改为：", string(ev.Kv.Value), "Revision：", ev.Kv.CreateRevision, ev.Kv.ModRevision)
				case mvccpb.DELETE:
					fmt.Println("删除了：", "Revision", ev.Kv.ModRevision)
				}
			}
		}
	}

}
