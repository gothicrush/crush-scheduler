package worker

import (
	"fmt"
	"github.com/gothicrush/crush-scheduler/common"
	"time"
)

// 任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent              // etcd 任务事件队列
	jobPlanTable      map[string]*common.JobSchedulePlan //任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo  // 任务执行表
	jobResultChan     chan *common.JobExecuteResult      // 任务结果队列
}

var (
	G_scheduler *Scheduler
)

// 处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		jobSchedulePlan, err := common.BuildJobSchedulePlan(jobEvent.Job)
		if err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE: //删除任务事件
		if _, exist := scheduler.jobPlanTable[jobEvent.Job.Name]; exist {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILLER: //强杀任务事件
		// 取消掉命令执行
		jobExecuteInfo, jobExecuting := scheduler.jobExecutingTable[jobEvent.Job.Name]

		if jobExecuting {
			jobExecuteInfo.CancelFunc() // 实现强杀
		}
	}
}

// 处理任务结果
func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	// 删除执行状态
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)

	// 生成执行日志
	if result.Err != common.ERR_LOCK_ALREADY_REQUIRED {
		jobLog := &common.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.ReadTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		}
		//存储日志到MongoDB
		G_logSink.Append(jobLog)
	}

	fmt.Println("任务执行完成：", result.ExecuteInfo.Job.Name, string(result.Output), result.Err)
}

// 调度协程
func (scheduler *Scheduler) scheduleLoop() {

	// 初始化(1s)
	scheduleAfter := scheduler.TrySchedule()

	// 调度的定时器
	scheduleTimer := time.NewTimer(scheduleAfter)

	// 定时任务
	for {
		select {
		case jobEvent := <-scheduler.jobEventChan: //监听任务变化事件
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: // 最近的任务过期了
		case jobResult := <-scheduler.jobResultChan: //监听任务执行结果
			scheduler.handleJobResult(jobResult)
		}
		// 调度最近一次任务
		scheduleAfter = scheduler.TrySchedule()
		scheduleTimer.Reset(scheduleAfter)
	}
}

// 尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	// 如果任务正在执行，跳过本次调度
	if _, executing := scheduler.jobExecutingTable[jobPlan.Job.Name]; executing {
		//fmt.Println("尚未跳过执行", jobPlan.Job.Name)
		return
	}

	// 构建执行状态信息
	jobExecuteInfo := common.BuildJobExecuteInfo(jobPlan)

	// 存入执行表中
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	// 执行任务
	fmt.Println("执行任务", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.ReadTime)
	G_executor.ExecuteJob(jobExecuteInfo)
}

// 重新计算任务调度状态
func (scheduler *Scheduler) TrySchedule() time.Duration {
	var scheduleAfter time.Duration

	// 如果任务表为空，则规定睡眠1s
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return scheduleAfter
	}

	// 获取当前时间
	var nearTime *time.Time
	now := time.Now()

	// 遍历所有任务
	for _, jobPlan := range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO: 尝试执行任务
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // 更新下次执行时间
		}

		//统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度间隔 = 最近要执行调度时间 - 当前时间
	scheduleAfter = (*nearTime).Sub(now)

	return scheduleAfter
}

// 推送任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

func InitScheduler() {
	// 赋值单例
	G_scheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}

	// 启动调度协程
	go G_scheduler.scheduleLoop()
}

// 回传任务执行结果
func (scheduler *Scheduler) PushJobResultJob(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
