package main

import (
	"fmt"
)

func gen(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func multiplyByTwo(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()
	return out
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	// если я все правильно понял и это классическая задача где данные передаются в каналы меж функциями
	in := gen(nums)
	out := multiplyByTwo(in)

	for v := range out {
		fmt.Println(v)
	}
}
