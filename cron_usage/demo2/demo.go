package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {

	var (
		cronJob *CronJob
		expr *cronexpr.Expression
		now time.Time
		scheduleTable map[string]*CronJob // key: 任务的名字,
	)
	scheduleTable = make(map[string]*CronJob)

	// 当前时间
	now = time.Now()


	// 定义两个 cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	scheduleTable["Job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册到调度表
	scheduleTable["Job2"] = cronJob

	// 需要有1个调度协程, 它定时检查所有的Cron任务, 谁过期了就执行谁
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now time.Time
		)
		// for循环 不能少， 检查任务调度表
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动协程，执行任务
					go func(jobName string) {
						fmt.Println("执行：",jobName)
					}(jobName)

					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println("下一次调度时间：", cronJob.nextTime)
				}
			}
		}

		select {
		case <- time.NewTimer(100*time.Millisecond).C:
		}
	}()
	time.Sleep(20* time.Second)
}
