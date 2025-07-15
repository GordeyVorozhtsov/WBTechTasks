package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	wg.Add(2)

	c := make(chan string)

	go func(chan string) {
		defer wg.Done()

		timer := time.After(time.Second * 3) // таймаут 3

		for {
			select {
			case c <- "oopps!":
				time.Sleep(time.Millisecond * 500)
			case <-timer:
				close(c)
				return
			}
		}
	}(c)

	go func(chan string) {
		defer wg.Done()

		for {
			val, ok := <-c
			if !ok {
				return
			}
			fmt.Println(val)
		}
	}(c)
	wg.Wait()

}
