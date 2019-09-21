package main

import (
	"fmt"
	"sort"
	"strings"
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
	wg.Wait()
}

func taskWrappper(task job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out) // обозначаем соседям, что задача завершилась
	task(in, out)
}

func SingleHash(in, out chan interface{}) {
	fmt.Printf("SingleHash:\n")
	for val := range in {
		data := fmt.Sprintf("%v", val)
		fmt.Printf("SingleHash: %s\n", data)
		out <- (DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data)))
	}
}

func MultiHash(in, out chan interface{}) {
	fmt.Printf("MultyHash:\n")
	const hashCount = 6
	for val := range in {
		data := fmt.Sprintf("%v", val)
		fmt.Printf("MultyHash: %s\n", data)
		tempResult := make([]string, 0, hashCount)
		for th := 0; th < hashCount; th++ {
			tempResult = append(tempResult, DataSignerCrc32(string(th)+data))
		}
		out <- strings.Join(tempResult, "")
	}
}

func CombineResults(in, out chan interface{}) {
	fmt.Printf("CombineResults:\n")
	dataPackage := make([]string, 0)
	for val := range in {
		data := fmt.Sprintf("%v", val)
		fmt.Printf("CombineResults got: %s\n", data)
		dataPackage = append(dataPackage, data)
	}
	fmt.Println("Start sorting")
	sort.Strings(dataPackage)
	fmt.Println("End sorting")
	out <- strings.Join(dataPackage, "_")
}
