package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// 单例
var (
	G_jobMgr *JobMgr
)

// 初始化管理器
func InitJobMgr() error {

	// etcd 连接配置
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,
		DialTimeout: time.Duration(G_config.EtcdTimeout) * time.Millisecond,
	}

	// 建立连接
	client, err := clientv3.New(config)

	if err != nil {
		return err
	}

	// 得到kv客户端
	kv := clientv3.NewKV(client)

	// 得到lease客户端
	lease := clientv3.NewLease(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return nil
}
