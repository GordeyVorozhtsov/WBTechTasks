package main

import (
	"fmt"
	"time"
)

// передаем сюда произвольное количество каналов с разным time.Duration и ожидаем ближайшего завершения
func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	// создаем под каждый канал горутину, которая будет ждать закрытия канала и закрывать out
	for _, c := range channels {
		go func(c <-chan interface{}) {
			<-c // пока канал открыт, чтение боликирует горутину -> читаем и можем закрыть наш out и завершить выполнение функции
			close(out)
		}(c)
	}
	return out

}
func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
