package main

import (
	"fmt"
	"time"
)

func sleep(duration int) {
	select { // блокировака как раз за счет селекта
	case <-time.After(time.Duration(duration) * time.Second):
		fmt.Println("конец таймеру")
	}
}
func main() {
	sleep(4)
	fmt.Println("продолжим!")
}
