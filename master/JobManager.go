package master

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gothicrush/crush-scheduler/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// 单例
var (
	G_jobManager *JobManager
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
	G_jobManager = &JobManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return nil
}

func (jobManager *JobManager) SaveJob(job *common.Job) (*common.Job, error) {
	//把任务保存到 /cron/jobs/任务名 -> json

	var oldJob common.Job

	// 键
	var jobKey = common.JOB_SAVE_DIR + job.Name

	// 值
	jobValue, err := json.Marshal(job)

	if err != nil {
		return nil, err
	}

	// 保存到 etcd
	putResponse, err := jobManager.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV())

	// 如果是更新，返回旧值
	if putResponse.PrevKv != nil {
		err = json.Unmarshal(putResponse.PrevKv.Value, &oldJob)

		if err != nil {
			return nil, nil
		}

		return &oldJob, nil
	}

	return nil, nil
}

func (jobManager *JobManager) DeleteJob(name string) (*common.Job, error) {

	// 键
	jobKey := common.JOB_SAVE_DIR + name

	// 从 etcd 中删除
	deleteResponse, err := jobManager.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())

	if err != nil {
		return nil, err
	}

	// 返回被删除的任务信息
	if len(deleteResponse.PrevKvs) == 0 {
		fmt.Println(jobKey)
		return nil, nil
	}

	var oldJob common.Job

	err = json.Unmarshal(deleteResponse.PrevKvs[0].Value, &oldJob)

	if err != nil {
		return nil, nil
	}

	return &oldJob, nil
}

func (jobManager *JobManager) ListJob() ([]*common.Job, error) {

	// 键
	dirKey := common.JOB_SAVE_DIR

	// 获取所有任务
	getReponse, err := jobManager.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix())

	if err != nil {
		return nil, err
	}

	var jobList []*common.Job

	for _, item := range getReponse.Kvs {
		var job common.Job = common.Job{}
		json.Unmarshal(item.Value, &job)

		jobList = append(jobList, &job)
	}

	return jobList, nil
}

func (jobManager *JobManager) KillJob(name string) error {

	// 更新 /cron/killer/任务名

	// 键
	killerKey := common.JOB_KILLER_DIR + name

	// 让worker监听一次put操作,使用租约，让其1秒后自动过期

	// 创建租约
	leaseResponse, err := jobManager.lease.Grant(context.TODO(), 1)

	if err != nil {
		return err
	}

	// 设置 killer
	_, err = jobManager.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseResponse.ID))

	if err != nil {
		return err
	}

	return nil
}
