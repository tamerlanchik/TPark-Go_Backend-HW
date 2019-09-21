package main

import (
	"fmt"
	"sync"
)

const (
	maxDataSize = 100
)

// сюда писать код
func main() {}

func ExecutePipeline(tasks ...job) {
	in := make(chan interface{}, maxDataSize)
	var out chan interface{}
	wg := &sync.WaitGroup{}

	for _, task := range tasks {
		out = make(chan interface{}, maxDataSize)
		wg.Add(1)
		go taskWrappper(task, in, out, wg)
		in = out
	}
	fmt.Println("Start waiting")
	wg.Wait()
	fmt.Println("End all tasks")
}

func taskWrappper(task job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out) // обозначаем соседям, что задача завершилась
	task(in, out)
}

func SingleHash(in, out chan interface{}) {
	//DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
