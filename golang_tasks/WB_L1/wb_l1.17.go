package main

import (
	"fmt"
)

func binarySearch(arr []int, x int) int {
	low := 0
	high := len(arr) - 1

	for low <= high {
		mid := (low + high) / 2

		if arr[mid] == x {
			return mid
		} else if arr[mid] < x {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func main() {
	fmt.Println(binarySearch([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10))
}
