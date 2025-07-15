package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// signal.Notify ждет ctrl+c, пишет в gracefulShutdown, закрывается канал done, отрабатывает выход из цикла в worker(),
// wg.Wait() ловит wg.Done() и программа завершается

var wg sync.WaitGroup

func main() {
	done := make(chan struct{})

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	workerPool(done, 5)

	sig := <-gracefulShutdown
	fmt.Println("got signal:", sig)

	close(done)
	wg.Wait()
}

func workerPool(done <-chan struct{}, n int) {
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(done, i)
	}
}

func worker(done <-chan struct{}, id int) {
	defer wg.Done()
	for {
		select {
		case <-done:
			fmt.Printf("got it!(worker№%d)\n", id)
			return
		default:
			fmt.Printf("worker %d waiting for ur ctrl+c\n", id)
			time.Sleep(time.Second)
		}
	}
}
