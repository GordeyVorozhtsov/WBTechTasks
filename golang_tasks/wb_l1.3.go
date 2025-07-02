package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	c := make(chan int, 1)
	// реализована бесконечная запись в канал(по условию)
	// выход из программы ctrl+c
	go func(chan int) {
		for {
			c <- 616
			time.Sleep(time.Second) // задержка в секунду для читабольности вывода
		}
	}(c)
	numOfWorker := 5 // количество воркеров
	workerpool(numOfWorker, c)
	wg.Wait()
}
func workerpool(n int, c chan int) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go workerRead(i, c)
	}
}

func workerRead(id int, c chan int) {
	defer wg.Done()
	for {
		num := <-c
		fmt.Printf("worker %d read num: %d\n", id, num)
	}
}
