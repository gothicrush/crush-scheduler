package _connect

import (
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {

	// etcd 客户端连接配置
	var etcdConfig clientv3.Config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	etcdClient, err := clientv3.New(etcdConfig)

	if err != nil {
		log.Println("clientv3.New failed : ", err)
		return
	}

	etcdClient = etcdClient
}
