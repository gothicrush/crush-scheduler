package worker

import (
	"context"
	"github.com/gothicrush/crush-scheduler/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

type JobManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

// 单例
var (
	G_jobManager *JobManager
)

// 初始化管理器
func InitJobManager() error {

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

	// 得到wather客户端
	watcher := clientv3.NewWatcher(client)

	// 赋值单例
	G_jobManager = &JobManager{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	// 启动任务监听
	G_jobManager.watchJobs()

	// 启动强杀监听
	G_jobManager.watchKiller()

	return nil
}

// 监听任务变化
func (jobManager *JobManager) watchJobs() error {
	// get /cron/jobs/ 目录下所有的任务
	getResponse, err := jobManager.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())

	if err != nil {
		return err
	}

	// 通过调度协程去执行任务
	for _, kvpair := range getResponse.Kvs {
		if job, err := common.UnpackJob(kvpair.Value); err == nil {
			// 构建一个 Save 事件
			jobEvent := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// 推送给调度协程
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	// 并用该 revision向后监听变化
	go func() {
		// 从 get 时刻的后续版本开始监听变化
		watchStartRevision := getResponse.Header.Revision + 1
		// 监听 /cron/jobs/ 目录后续变化
		watchChan := jobManager.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR,
			clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResponse := range watchChan {
			for _, watchEvent := range watchResponse.Events {
				var job *common.Job
				var jobEvent *common.JobEvent

				switch watchEvent.Type {
				case mvccpb.PUT: //任务保存事件
					//反序列化job，推送给调度协程
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构建一个 Save 事件，准备进行推送
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE: //任务删除事件
					jobName := common.ExtractJobName(string(watchEvent.Kv.Key))
					job := &common.Job{Name: jobName}
					// 构建一个 Delete 事件，准备进行推送
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}

				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()

	return nil
}

// 创建任务执行锁
func (jobManager *JobManager) CreateJobLock(jobName string) *JobLock {
	// 返回一把锁
	return InitJobLock(jobName, jobManager.kv, jobManager.lease)
}

// 监听强杀任务
func (jobManager *JobManager) watchKiller() {
	go func() {
		// 监听 /cron/jobs/ 目录后续变化
		watchChan := jobManager.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR,
			clientv3.WithPrefix())
		// 处理监听事件
		for watchResponse := range watchChan {
			for _, watchEvent := range watchResponse.Events {
				//var job *common.Job
				//var jobEvent *common.JobEvent

				switch watchEvent.Type {
				case mvccpb.PUT: //杀死任务事件
					jobName := common.ExtractKillerName(string(watchEvent.Kv.Key))
					job := &common.Job{Name: jobName}
					jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILLER, job)
					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: //killer标记过期，被自动删除，不关心
				}
			}
		}
	}()
}
