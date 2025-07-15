package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	arr := []int{2, 4, 6, 8, 10}
	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		go square(arr[i])
	}
	wg.Wait()
}
func square(i int) {
	defer wg.Done()
	fmt.Printf("square of %d is %d\n", i, i*i)
}
