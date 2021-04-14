package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	var (
		expr     *cronexpr.Expression
		err      error
		now      time.Time
		nextTime time.Time
	)

	// linux crontab
	// 支持 秒 粒度，年配置，（枚举到2099）
	// 分钟（0-59）、小时（0-23）、哪天（1-31）、哪月（1-12）、星期几（0-6）
	// 已经声明的变量，不要用:=
	if expr, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	now = time.Now()
	nextTime = expr.Next(now)

	fmt.Println(now, nextTime)

	// 等待定时器超时
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了：", nextTime)
	})
	
	time.Sleep(time.Second * 5)
}
