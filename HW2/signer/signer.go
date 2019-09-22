package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maxDataSize = 100
)

type pack struct {
	in  chan interface{}
	out chan interface{}
}

func (p *pack) Init(len1, len2 int) {
	p.in = make(chan interface{}, len1)
	p.out = make(chan interface{}, len2)
}

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

//---------------------------------------------
type HashExecutor struct {
	Channel chan pack
}

func (h *HashExecutor) PipelineCrc32() {
	wg := &sync.WaitGroup{}
	fmt.Println("Start PipelineCrs32")
	for task := range h.Channel {
		//fmt.Println("PipelineCrc32: value")
		data := fmt.Sprintf("%v", <-task.in)
		fmt.Printf("PipelineCrc32-value: %s\n", data)
		wg.Add(1)
		go func(f func(data string) string, data string, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			//defer close(out)
			out <- f(data)
			fmt.Printf("PipelineCrc32-inside: done - %s\n", data)
		}(DataSignerCrc32, data, task.out, wg)
	}
	wg.Wait()
}

func (h *HashExecutor) PipelineMd5() {
	fmt.Println("Start PipelineMD5")
	for task := range h.Channel {
		data := fmt.Sprintf("%v", <-task.in)
		fmt.Printf("PipelineMd5-value: %s\n", data)
		task.out <- DataSignerMd5(data)
		fmt.Printf("PipelineMd5: step done - %s\n", data)
	}
}
func (h *HashExecutor) Init() {
	h.Channel = make(chan pack, maxDataSize*3)
}
func (h *HashExecutor) AddPack(data pack) {
	h.Channel <- data
}

var Crc32Executor = HashExecutor{}
var Md5Executor = HashExecutor{}

func SingleHash(in, out chan interface{}) {
	timer := time.Tick(time.Second)
	go func() {
		for t := range timer {
			fmt.Println("\nTIME: ", t)
		}
	}()
	fmt.Printf("SingleHash:\n")
	Crc32Executor.Init()
	Md5Executor.Init()
	go Crc32Executor.PipelineCrc32()
	go Md5Executor.PipelineMd5()

	wg := &sync.WaitGroup{}
	for val := range in {
		fmt.Printf("SingleHash: %s\n", fmt.Sprintf("%v", val))
		//out <- (DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data)))
		wg.Add(1)
		go func(val interface{}) {
			defer wg.Done()
			pack1 := pack{}
			pack1.Init(1, 1)
			pack1.in <- val

			pack2 := pack{}
			pack2.Init(1, 1)
			pack2.in <- val

			pack3 := pack{}
			pack3.Init(1, 1)

			Crc32Executor.AddPack(pack1)
			Md5Executor.AddPack(pack2)
			var s1, s2 string
			for i := 0; i < 3; i++ {
				select {
				case data := <-pack1.out:
					s1 = data.(string) + "~"
					fmt.Printf("SingleHash-inside-%s: got pack1.out: %s\n", data, s1)
				case data := <-pack2.out:
					pack3.in <- data
					fmt.Printf("SingleHash-inside-%s: got pack2.out: %s\n", data, s1)
					Crc32Executor.AddPack(pack3)
				case data := <-pack3.out:
					s2 = data.(string)
					fmt.Printf("SingleHash-inside-%s: got pack3.out: %s\n", data, s1)
				}
			}
			out <- s1 + s2
			fmt.Printf("Single-hash< result\n")
		}(val)
		//time.Sleep(10 * time.Second)

	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	time.Sleep(1 * time.Second)
	fmt.Printf("\nMultyHash:\n\n")
	const hashCount = 6
	wgSuper := &sync.WaitGroup{}
	for val := range in {
		wgSuper.Add(1)
		go func(val interface{}) {
			defer wgSuper.Done()
			wg := &sync.WaitGroup{}
			data := val.(string)
			fmt.Printf("MultyHash: %s\n", data)
			tempResult := make([]string, hashCount, hashCount)
			for th := 0; th < hashCount; th++ {
				wg.Add(1)
				go func(data string, th int) {
					defer wg.Done()
					pack1 := pack{}
					pack1.Init(1, 1)
					pack1.in <- (strconv.Itoa(th) + data)
					Crc32Executor.AddPack(pack1)
					res := (<-pack1.out).(string)
					tempResult[th] = res
					fmt.Printf("MultiHash-inside (%s): done %s\n", data, res)
				}(data, th)
			}
			wg.Wait()
			result := strings.Join(tempResult, "")
			fmt.Printf("\nMultiHash-result: %s\n\n", result)
			out <- result
		}(val)
		//time.Sleep(10 * time.Second)
	}
	wgSuper.Wait()
}

func CombineResults(in, out chan interface{}) {
	time.Sleep(4 * time.Second)
	fmt.Printf("CombineResults:\n")
	dataPackage := make([]string, 0)
	for val := range in {
		//data := fmt.Sprintf("%v", val)
		data := val.(string)
		fmt.Printf("CombineResults got: %s\n", data)
		dataPackage = append(dataPackage, data)
	}
	fmt.Println("Start sorting")
	sort.Strings(dataPackage)
	fmt.Println("End sorting")
	out <- strings.Join(dataPackage, "_")
}
