package main

func main() {
	ch := make(chan int) //обявляем канал обязательно через make
	go func() {          //пишем в канал
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch) //надо закрыть канал чтобы не было бесконечного чтения и дедлока
	}()
	for n := range ch { //читаем из канала
		println(n)
	}
}
