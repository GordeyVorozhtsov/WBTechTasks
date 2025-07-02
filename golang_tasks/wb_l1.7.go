package main

import (
	"fmt"
	"sync"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
	sm sync.Map
)

func main() {
	byMutex()
	bySyncMap()
}

func byMutex() {
	m := make(map[int]int)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mu.Lock()
			m[i] = i
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	for key, value := range m {
		fmt.Println(key, value) // вывод будет рандомный потому что в мапе нет порядка
	}
}

func bySyncMap() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			sm.Store(i, i)
		}
	}()
	wg.Wait()

	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("%v: %v\n", key, value)
		return true
	})
}
