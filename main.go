package main

import (
	"fmt"
	"icrawler/task"
	"time"
)

func main() {
	start := time.Now()
	task.Run()
	//task.BilibiliPage()
	duration := time.Since(start)
	fmt.Printf("总耗时：%s\n\n", duration)
}
