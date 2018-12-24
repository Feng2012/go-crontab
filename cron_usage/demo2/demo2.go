package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

//CronJob任务
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time // expr.Next(now)
}

func main() {
	// 需要一个调用协程，定时检查所有的cron任务，谁过期就执行谁
	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob // 调度表
	)
	scheduleTable = make(map[string]*CronJob)
	// 当前时间
	now = time.Now()
	// 1 定义两个CronJob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册到调度表
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/7 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	scheduleTable["job2"] = cronJob

	// 启动调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)

		// 定时检查任务调度表
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				// 下次调度时间小于等于当前时间，说明任务到期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程，执行这个任务
					go func(jobName string) {
						fmt.Println("执行:", jobName)
					}(jobName)

					// 计算下次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, " 下次执行时间: ", cronJob.nextTime)
				}
			}
			// 睡眠100ms
			select {
			case <-time.NewTimer(100 * time.Millisecond).C: // 100ms后可读，返回
			}

			// time.Sleep(100 * time.Millisecond) // 睡眠100ms
		}
	}()

	time.Sleep(100 * time.Second)
}
