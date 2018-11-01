package worker

import (
	"github.com/gothicrush/crush-scheduler/common"
	"math/rand"
	"os/exec"
	"time"
)

// 执行器
type Executor struct {
}

var (
	G_executor *Executor
)

// 执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	// 开启一个新协程执行shell命令
	go func() {
		// 获取命令对象
		cmd := exec.CommandContext(info.CancelCtx, "C:\\cygwin\\bin\\bash.exe", "-c", info.Job.Command)

		// 创建任务执行结果，包括执行信息与执行结果
		result := &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		// 创建分布式锁
		jobLock := G_jobManager.CreateJobLock(info.Job.Name)

		// 尝试上锁
		// 随机睡眠，避免抢锁倾斜
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		// 尝试上锁
		err := jobLock.TryLock()
		// 延时释放锁
		defer jobLock.Unlock()

		if err != nil { // 上锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else { // 上锁成功
			// 任务开始时间
			result.StartTime = time.Now()

			// 执行命令并捕获输出
			output, err := cmd.CombinedOutput()

			// 任务执行结果
			result.Output = output
			result.Err = err

			// 任务结束时间
			result.EndTime = time.Now()

		}

		// 将执行结果返回给Scheduler，Scheduler会从executingTable中删除执行记录
		G_scheduler.PushJobResultJob(result)
	}()
}

func InitExecutor() error {

	G_executor = &Executor{}

	return nil
}
