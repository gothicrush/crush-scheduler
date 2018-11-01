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

// 初始化任务管理器
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
	// 获取当前已经存在的所有的任务
	getResponse, err := jobManager.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix())

	if err != nil {
		return err
	}

	// 通过调度协程去执行已经存在的任务
	for _, kvpair := range getResponse.Kvs {
		// 将值转为任务
		if job, err := common.UnpackJob(kvpair.Value); err == nil {
			// 每个派发到执行器的任务都要封装为一个任务事件
			// 构建一个 Save 事件
			jobEvent := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// 推送给调度协程
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	// 开启协程，用该revision向后监听后来新增的任务
	go func() {
		// 从 get 时刻的后续版本开始监听变化
		watchStartRevision := getResponse.Header.Revision + 1
		// 监听 /cron/jobs/ 目录后续变化
		watchChan := jobManager.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR,
			clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResponse := range watchChan {
			for _, watchEvent := range watchResponse.Events {
				// 从监听事件中获取任务，并封装为任务事件，发给执行器
				var job *common.Job
				var jobEvent *common.JobEvent

				switch watchEvent.Type {
				case mvccpb.PUT: //任务保存事件
					//反序列化job，推送给调度协程
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue // 失败则不管了
					}
					// 构建一个 Save 事件，准备进行推送
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE: //任务删除事件
					//提取出需要删除任务的名字
					jobName := common.ExtractJobName(string(watchEvent.Kv.Key))
					//封装为一个任务
					job := &common.Job{Name: jobName}
					// 构建一个 Delete 事件，准备进行推送
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}

				// 将任务事件推送给执行器
				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()

	return nil
}

// 创建分布式锁，并返回
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
			// 处理监听事件
			for _, watchEvent := range watchResponse.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //杀死任务事件
					// 获取强杀任务名
					jobName := common.ExtractKillerName(string(watchEvent.Kv.Key))
					// 封装为一个强杀任务
					job := &common.Job{Name: jobName}
					// 封装为一个强杀事件
					jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILLER, job)
					// 将强杀事件推送给执行器执行
					G_scheduler.PushJobEvent(jobEvent)
				}
			}
		}
	}()
}
