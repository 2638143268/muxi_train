package main

import (
	"fmt"
	"sync"
)

// 1.并发加数
// 要求开启50个协程去并发地将一个变量从0增加到5000，也就是说每个协程需要加给这个数加100，但是
// 每个协程每次只能给这个变量加1

var (
	mu sync.Mutex
)

func add(a *int, wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		mu.Lock()
		*a++
		mu.Unlock()
	}
	wg.Done()
}
func main() {
	var wg sync.WaitGroup
	var a int = 0
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go add(&a, &wg)
	}
	wg.Wait()
	fmt.Println(a)
}
