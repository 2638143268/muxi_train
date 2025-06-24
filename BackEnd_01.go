package main

import "fmt"

func Prime(n int) []int {
	var arr []int = []int{}
	for i := 2; i <= n; i++ {
		var set bool = true
		for j := 2; j < i; j++ {
			if i%j == 0 {
				set = false
				break
			}
		}
		if set {
			arr = append(arr, i)
		}
	}
	return arr
}

func main() {
	var n int = 100
	var arr []int = Prime(n)
	fmt.Println(arr)
}
