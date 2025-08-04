package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 12}
	fmt.Println(rotate(a, 3))
}
func rotate(nums []int, k int) []int {
	j := len(nums) - 1
	for i := j - k; i >= 0; i-- {
		nums[i], nums[j] = nums[j], nums[i]
		j--
	}
	return nums
}
