package main

import (
	"fmt"
	"time"

	"localhost.com/my_web/logic/config"
	"localhost.com/my_web/logic/excel"
)

func main() {
	fmt.Println("执行涨停数据收集", time.Now().Format("2006-01-02 15:04:05"))
	// 这里调用你的任务函数
	excel.SetStockExcel(true)
	time.Sleep(time.Second * 3)
	config.SetConfigFile(true)
	fmt.Println("执行涨停数据收集完毕", time.Now().Format("2006-01-02 15:04:05"))
}
