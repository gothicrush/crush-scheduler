package master

import (
	"context"
	"github.com/gothicrush/crush-scheduler/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type WorkerManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	// 单例
	G_workerManager *WorkerManager
)

// 列举所有的worker
func (workerManager *WorkerManager) ListWorkers() ([]string, error) {

	var workerArr []string

	getResp, err := workerManager.kv.Get(context.TODO(), common.JOB_WORKER_DIR, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	//解析每个节点的IP
	for _, kv := range getResp.Kvs {
		workerIP := common.ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}

	return workerArr, nil
}

func InitWorkerManager() error {

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

	G_workerManager = &WorkerManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return nil
}
