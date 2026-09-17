package main

import "fmt"

func twoSum(nums []int, target int) []int {
	positions := map[int]int{}
	for i := 0; i < len(nums); i++ {
		positions[nums[i]] = i
	}
	for i := 0; i < len(nums); i++ {
		if positions[target-nums[i]] != 0 && positions[target-nums[i]] != i {
			return []int{i, positions[target-nums[i]]}
		}
	}
	return nil
}

func main() {
	fmt.Println(twoSum([]int{2, 5, 5, 11}, 10))
}
