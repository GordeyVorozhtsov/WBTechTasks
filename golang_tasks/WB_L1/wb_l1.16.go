package main

import (
	"fmt"
)

func quickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	pivotIndex := len(arr) / 2
	pivotValue := arr[pivotIndex]
	var left, equal, right []int
	for _, v := range arr {
		if v < pivotValue {
			left = append(left, v)
		} else if v == pivotValue {
			equal = append(equal, v)
		} else {
			right = append(right, v)
		}
	}
	left = quickSort(left)
	right = quickSort(right)
	return append(append(left, equal...), right...)
}

func main() {
	fmt.Println(quickSort([]int{5, 1, 7, 2, 4, 9, 3, 8, 6, 10}))
}
