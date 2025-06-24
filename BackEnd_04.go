package main

import "fmt"

func merge(nums1 []int, m int, nums2 []int, n int) {
	i, j, k := m-1, n-1, m+n-1
	for i >= 0 && j >= 0 {
		if nums1[i] > nums2[j] {
			nums1[k] = nums1[i]
			i--
		} else {
			nums1[k] = nums2[j]
			j--
		}
		k--
	}
	for j >= 0 {
		nums1[k] = nums2[j]
		j--
		k--
	}
}

func main() {
	nums1 := make([]int, 6)
	copy(nums1, []int{1, 3, 5})
	nums2 := []int{2, 4, 6}
	merge(nums1, 3, nums2, 3)
	fmt.Println(nums1)
}
