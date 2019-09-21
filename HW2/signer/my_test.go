package main

// Тестируем однократный обмен данными
/*
func TestOne(t *testing.T) {
	numb := 123
	var newNumb int
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			time.Sleep(time.Millisecond)
			out <- numb * 3
		}),
		job(func(in, out chan interface{}) {
			value := <-in
			time.Sleep(time.Millisecond)
			out <- value
			newNumb = value.(int)
		}),
	}

	ExecutePipeline(freeFlowJobs...)

	if numb*3 != newNumb {
		t.Errorf("Wrong value got: %d instead %d", newNumb, numb*3)
	}
}
*/
// Тест на продолжительный поток данных
/*
func TestTwo(t *testing.T) {
	var result1, result2 int
	const (
		count      = maxDataSize
		maxTimeout = 30
	)
	rand.Seed(1234)
	freeFlowJobs := []job{
		job(func(in, out chan interface{}) {
			for i := 0; i < count; i++ {
				out <- i
				timeout := rand.Intn(maxTimeout)
				time.Sleep(time.Duration(timeout) * time.Millisecond)
			}
		}),
		job(func(in, out chan interface{}) {
			sum := 0
			for val := range in {
				sum += val.(int)
				out <- val.(int)
				timeout := rand.Intn(maxTimeout)
				time.Sleep(3 * time.Duration(timeout) * time.Millisecond)
			}
			result1 = sum
		}),
		job(func(in, out chan interface{}) {
			sum := 0
			for val := range in {
				sum += val.(int)
				out <- val.(int)
			}
			result2 = sum
		}),
	}
	ExecutePipeline(freeFlowJobs...)

	trueResult := count * (count - 1) / 2
	if result1 != trueResult {
		t.Errorf("Wrong1 : got %d instead %d", result1, trueResult)
	}
	if result2 != trueResult {
		t.Errorf("Wrong2 : got %d instead %d", result2, trueResult)
	}
}
*/
