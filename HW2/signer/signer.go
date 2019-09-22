package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	maxDataSize = 100
)

type ChanPack struct {
	in  chan interface{}
	out chan interface{}
}

func (p *ChanPack) Init(len1, len2 int) {
	p.in = make(chan interface{}, len1)
	p.out = make(chan interface{}, len2)
}

// сюда писать код

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

//---------------------------------------------

type HashExecutor struct {
	Channel chan ChanPack
}

func (h *HashExecutor) PipelineCrc32() {
	wg := &sync.WaitGroup{}
	for task := range h.Channel {
		wg.Add(1)
		go func(f func(data string) string, task ChanPack) {
			defer wg.Done()
			data := fmt.Sprintf("%v", <-task.in)
			task.out <- f(data)
		}(DataSignerCrc32, task)
	}
	wg.Wait()
}

func (h *HashExecutor) PipelineMd5() {
	for task := range h.Channel {
		data := fmt.Sprintf("%v", <-task.in)
		task.out <- DataSignerMd5(data)
	}
}
func (h *HashExecutor) Init() {
	h.Channel = make(chan ChanPack, maxDataSize)
}
func (h *HashExecutor) AddTask(data ChanPack) {
	h.Channel <- data
}

var Crc32Executor = HashExecutor{}
var Md5Executor = HashExecutor{}

func SingleHash(in, out chan interface{}) {
	Crc32Executor.Init()
	Md5Executor.Init()
	go Crc32Executor.PipelineCrc32()
	go Md5Executor.PipelineMd5()

	wg := &sync.WaitGroup{}
	for val := range in {
		wg.Add(1)
		go func(val interface{}) {
			defer wg.Done()
			crc32Pack := ChanPack{}
			crc32Pack.Init(1, 1)
			crc32Pack.in <- val

			md5Pack := ChanPack{}
			md5Pack.Init(1, 1)
			md5Pack.in <- val

			md5_to_crc32Pack := ChanPack{}
			md5_to_crc32Pack.Init(1, 1)

			Crc32Executor.AddTask(crc32Pack)
			Md5Executor.AddTask(md5Pack)
			var s1, s2 string
			for i := 0; i < 3; i++ {
				select {
				case data := <-crc32Pack.out:
					s1 = data.(string) + "~"
				case data := <-md5Pack.out:
					md5_to_crc32Pack.in <- data
					Crc32Executor.AddTask(md5_to_crc32Pack)
				case data := <-md5_to_crc32Pack.out:
					s2 = data.(string)
				}
			}
			out <- s1 + s2
		}(val)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	const hashCount = 6
	wgSuper := &sync.WaitGroup{}

	for val := range in {
		wgSuper.Add(1)

		go func(val interface{}) {
			defer wgSuper.Done()

			wg := &sync.WaitGroup{}
			tempResult := make([]string, hashCount, hashCount)
			data := val.(string)

			for th := 0; th < hashCount; th++ {
				wg.Add(1)
				go func(data string, th int) {
					defer wg.Done()

					crc32Pack := ChanPack{}
					crc32Pack.Init(1, 1)
					crc32Pack.in <- (strconv.Itoa(th) + data)

					Crc32Executor.AddTask(crc32Pack)

					tempResult[th] = (<-crc32Pack.out).(string)
				}(data, th)
			}
			wg.Wait()
			out <- strings.Join(tempResult, "")
		}(val)
	}
	wgSuper.Wait()
}

func CombineResults(in, out chan interface{}) {
	dataPackage := make([]string, 0)
	for val := range in {
		dataPackage = append(dataPackage, val.(string))
	}
	sort.Strings(dataPackage)
	out <- strings.Join(dataPackage, "_")
}
