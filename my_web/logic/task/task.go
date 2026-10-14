package task

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"localhost.com/my_web/logic/config"
	"localhost.com/my_web/logic/excel"
	"localhost.com/my_web/logic/task/bidtask"
)

/*
	处理一些定时任务，如：9.20加载集合竞价数据，15.30处理昨日二板数据等
*/

func AutoTask(ctx context.Context) {
	fmt.Println("程序启动:", time.Now().Format("2006-01-02 15:04:05"))

	c := cron.New(cron.WithLocation(time.Local))
	defer c.Stop() // 确保在函数退出时停止cron

	// 使用通道来通知任务完成
	done := make(chan struct{}, 1)

	// 添加任务1：9:20执行
	_, err := c.AddFunc("20 9 * * *", func() {
		fmt.Println("执行集合竞价数据收集:", time.Now().Format("2006-01-02 15:04:05"))
		bidtask.GethqInfo()
		// done <- struct{}{} // 通知任务完成
		fmt.Println("执行集合竞价数据收集完毕:", time.Now().Format("2006-01-02 15:04:05"))
	})

	if err != nil {
		fmt.Printf("添加集合竞价数据收集任务失败: %v\n", err)
		return
	}

	// 添加任务2：15:30执行
	_, err2 := c.AddFunc("30 15 * * *", func() {
		fmt.Println("执行涨停数据收集", time.Now().Format("2006-01-02 15:04:05"))
		// 这里调用你的任务函数
		excel.SetStockExcel(true)
		time.Sleep(time.Second * 3)
		config.SetConfigFile(true)
		fmt.Println("执行涨停数据收集完毕", time.Now().Format("2006-01-02 15:04:05"))
		done <- struct{}{} // 通知任务完成
	})

	if err2 != nil {
		fmt.Printf("添加涨停数据收集任务失败: %v\n", err2)
		return
	}

	// 启动定时任务
	c.Start()

	// 等待任一任务完成或上下文取消
	select {
	case <-done:
		fmt.Println("任务执行完毕，程序退出")
	case <-ctx.Done():
		fmt.Println("上下文已取消，程序退出")
	}
}
