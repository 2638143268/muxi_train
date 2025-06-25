package main

import (
	"fmt"
	"sync"
	"time"
)

// 2.交替打印
// 开启两个协程，来交替打印字母(A-Z)和数字(1-26)，最终呈现的效果大概是
// A1B2C3D4....

var (
	wg  sync.WaitGroup
	ch1 = make(chan int, 1)
	ch2 = make(chan int, 1)
)

func printAZ() {
	for i := 0; i < 26; i++ {
		a := <-ch2
		for a == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Print(string('A' + i))
		ch1 <- 1
	}
	wg.Done()
}
func print126() {
	for i := 0; i < 26; i++ {
		b := <-ch1
		for b == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Print(i + 1)
		ch2 <- 1
	}
	wg.Done()
}
func main() {
	ch2 <- 1
	wg.Add(1)
	go printAZ()
	wg.Add(1)
	go print126()
	wg.Wait()
}
