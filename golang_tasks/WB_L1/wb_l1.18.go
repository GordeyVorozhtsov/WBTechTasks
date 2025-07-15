package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var wg sync.WaitGroup

type Inc struct {
	Counter int
}

func (a *Inc) someWork(x int) {
	defer wg.Done()
	for i := 0; i < x; i++ {

		mu.Lock()
		a.Counter++
		mu.Unlock()
	}
}
func main() {
	inc := &Inc{}

	wg.Add(3)

	go inc.someWork(33)
	go inc.someWork(27)
	go inc.someWork(40)

	wg.Wait()

	x := inc.Counter
	fmt.Println(x)
}
