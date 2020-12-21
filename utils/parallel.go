package utils

import (
	"sync"
)

// Config 配置并行任务
//	Data: 任何执行完毕后返回的数据
//	Numbers: 并行的数量
//	WaitGroup: 等待 goroutine 执行完毕 Add / Done / Wait
type Config struct {
	Data      []interface{}
	Numbers   chan interface{}
	WaitGroup sync.WaitGroup
}

// MultiRun 执行并行任务
//	run 并行任务中实际执行任务的函数
//	taskList 等待处理的任务
//	thread 并行的数量
func MultiRun(run func(string) interface{}, taskList []string, thread int) (data []interface{}) {
	taskNum := len(taskList)

	config := new(Config)
	config.Numbers = make(chan interface{}, thread)
	config.WaitGroup.Add(taskNum)

	for i := 0; i < taskNum; i++ {
		task := taskList[i]
		config.Numbers <- task
		go func() {
			config.Data = append(config.Data, run(task))
			config.WaitGroup.Done()
			<-config.Numbers
		}()
	}
	config.WaitGroup.Wait()

	data = make([]interface{}, 0)
	data = config.Data
	return
}
