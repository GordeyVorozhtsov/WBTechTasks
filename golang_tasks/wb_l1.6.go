package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// выход по условию
	exitByCondition()

	// выход через канал
	exitByChannel()

	// выход с контекстом
	exitByContext()

	// выход с runtime.Goexit()
	exitByGoexit()
}

func exitByCondition() {
	stop := false
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if stop { // если будет stop = true тогда программа дойдет до return
				fmt.Println("got condition true")
				return
			}
			fmt.Println("waiting for chanching condition")
			time.Sleep(time.Millisecond * 500)
		}
	}()

	time.Sleep(time.Second)
	stop = true
	wg.Wait()
}

func exitByChannel() {
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopCh:
				fmt.Println("got closing")
				return
			default:
				fmt.Println("waiting for closing channel")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()

	time.Sleep(time.Second)
	close(stopCh)
	wg.Wait()
}

func exitByContext() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("got ctx cancel")
				return
			default:
				fmt.Println("waiting for cancel ctx")
				time.Sleep(time.Millisecond * 500)
			}
		}
	}(ctx)

	time.Sleep(time.Second)
	cancel()
	wg.Wait()
}

func exitByGoexit() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("starting")
		time.Sleep(time.Second)
		fmt.Println("calling runtime.Goexit()")
		runtime.Goexit()
		fmt.Println("ist wouldnt be printing")
	}()
	wg.Wait()
}
