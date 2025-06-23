package main

import "fmt"

func Deduplicate(nums []int) []int {
	i := 0
	for i < len(nums) {
		j := i + 1
		for j < len(nums) {
			if nums[i] == nums[j] {
				nums = append(nums[:j], nums[j+1:]...)
			} else {
				j++
			}
		}
		i++
	}
	return nums
}

func main() {
	arr := []int{1, 2, 3, 2, 5, 3}
	arr = Deduplicate(arr)
	fmt.Println(arr)
}
