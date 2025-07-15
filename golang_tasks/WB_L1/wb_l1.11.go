package main

import (
	"fmt"
)

func main() {
	fmt.Println(findCrossing([]int{1, 2, 3}, []int{2, 3, 4}))
	fmt.Println(findCrossing([]int{777, 778, 779, 780, 781, 782}, []int{779, 780, 781, 782, 783, 784, 785}))

}
func findCrossing(a, b []int) []int {
	res := []int{}
	m := make(map[int]bool)
	for _, val := range a {
		m[val] = true
	}
	for _, val := range b {
		if m[val] {
			res = append(res, val)
		}
	}
	return res
}
