package worker

import (
	"context"
	"github.com/gothicrush/crush-scheduler/common"
	"go.etcd.io/etcd/clientv3"
)

// 分布式锁
type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string             // 任务名，也是锁名
	cancelFunc context.CancelFunc //用于终止自动续租
	leaseID    clientv3.LeaseID   //锁的租约ID，用于释放锁过程中立即释放租约
	isLocked   bool               //是否已经上锁
}

// 创建一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) *JobLock {
	jobLock := &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}

	return jobLock
}

// 尝试上锁
func (jobLock *JobLock) TryLock() error {

	// 创建5秒租约
	leaseResp, err := jobLock.lease.Grant(context.TODO(), 5)

	if err != nil {
		return err
	}

	// 自动续租
	// 创建可以取消的context
	cancleContext, cancelFunc := context.WithCancel(context.TODO())
	leaseID := leaseResp.ID
	keepRespChan, err := jobLock.lease.KeepAlive(cancleContext, leaseID)

	if err != nil {
		cancelFunc()                                  // 取消自动续租
		jobLock.lease.Revoke(context.TODO(), leaseID) // 立即释放租约
		return err
	}

	// 处理续租应答的协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	// 创建事务txn
	txn := jobLock.kv.Txn(context.TODO())

	// 锁路径
	lockKey := common.JOB_LOCK_DIR + jobLock.jobName

	// 事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet(lockKey))

	// 提交事务
	var txnResp *clientv3.TxnResponse
	if txnResp, err = txn.Commit(); err != nil {
		cancelFunc()                                  // 取消自动续租
		jobLock.lease.Revoke(context.TODO(), leaseID) // 立即释放租约
		return err
	}

	// 成功则返回，失败释放租约
	if !txnResp.Succeeded {
		cancelFunc()                                  // 取消自动续租
		jobLock.lease.Revoke(context.TODO(), leaseID) // 立即释放租约
		return common.ERR_LOCK_ALREADY_REQUIRED
	}

	// 抢锁成功
	jobLock.leaseID = leaseID
	jobLock.cancelFunc = cancelFunc
	jobLock.isLocked = true

	return nil
}

// 释放锁
func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
		jobLock.cancelFunc()                                  //取消自动续租协程
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseID) // 释放租约
	}
}
