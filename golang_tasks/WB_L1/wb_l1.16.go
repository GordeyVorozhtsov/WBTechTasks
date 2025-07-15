package main

import (
	"fmt"
)

// первый алгоритм что мне пришел в голову
func quickSort(a []int) []int {
	zeroIdx := 0
	for zeroIdx < len(a) {
		for idx := 0; idx < len(a)-1; idx++ {
			zeroIdx++
			if a[idx] > a[idx+1] {
				a[idx], a[idx+1] = a[idx+1], a[idx]
				zeroIdx = 0
			}
		}
	}
	return a
}
func main() {
	fmt.Println(quickSort([]int{5, 1, 7, 2, 4, 9, 3, 8, 6, 10}))
}
