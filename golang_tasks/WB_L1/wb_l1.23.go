package main

import "fmt"

func deleteIndex(slice []int, i int) []int {
	return append(slice[:i], slice[i+1:]...)
}

func main() {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	arr = deleteIndex(arr, 4)
	fmt.Println(arr)
}
