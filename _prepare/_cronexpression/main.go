package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

/*
启动一个调度协程，定时检查所有的Cron任务，谁过期了就执行谁

cron表达式计算出的执行时间与运行cron表达式的时间没有关系
time.AfterFunc(d Duration, f func(){})
*/

type CronJob struct {
	expression *cronexpr.Expression
	nextTime   time.Time
}

func main() {

	now := time.Now()

	//调度表
	var scheduleTable map[string]*CronJob = make(map[string]*CronJob)

	// 1. 第一个 cronjob
	var expression1 *cronexpr.Expression = cronexpr.MustParse("*/5 * * * * * *")

	cj1 := &CronJob{
		expression: expression1,
		nextTime:   expression1.Next(now),
	}

	scheduleTable["cj1"] = cj1

	// 2. 第二个 cronjob
	var expression2 *cronexpr.Expression = cronexpr.MustParse("*/7 * * * * * *")

	cj2 := &CronJob{
		expression: expression2,
		nextTime:   expression2.Next(now),
	}

	scheduleTable["cj2"] = cj2

	for {
		now := time.Now()
		for jobName, cronJob := range scheduleTable {
			if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
				//启动一个协程，执行这个任务
				go func(jobName string) {
					fmt.Println("执行：", jobName, time.Now())
				}(jobName)

				//计算下一次调度时间
				cronJob.nextTime = cronJob.expression.Next(now)
				fmt.Println(jobName, "下次执行时间:", cronJob.nextTime)
			}
		}

		// 睡眠100毫秒
		select {
		// 将在100毫秒后可读，返回
		case <-time.NewTimer(100 * time.Millisecond).C:
		}

		// time.Sleep(100 * time.Millisecond)
	}
}
