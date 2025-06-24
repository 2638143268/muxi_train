package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func BuildTree(nums []int) *TreeNode {
	if len(nums) == 0 || nums[0] == -1 {
		return nil
	}
	root := &TreeNode{Val: nums[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(nums) {
		node := queue[0]
		queue = queue[1:]
		if i < len(nums) && nums[i] != -1 {
			node.Left = &TreeNode{Val: nums[i]}
			queue = append(queue, node.Left)
		}
		i++
		if i < len(nums) && nums[i] != -1 {
			node.Right = &TreeNode{Val: nums[i]}
			queue = append(queue, node.Right)
		}
		i++
	}
	return root
}
func printTree(root *TreeNode) {
	if root == nil {
		fmt.Println("Empty tree")
		return
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			fmt.Print(node.Val, " ")
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		fmt.Println()
	}
}
func main() {
	nums := []int{1, 2, 3, 9, 8, 4, 5}
	root := BuildTree(nums)
	printTree(root)
}
