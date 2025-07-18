package main

import "fmt"

// здесь мы обьявляем переменную x = 1 и тк defer выполняется последним,
// он инкрементирует переменную х (х = 2)
func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}

// здесь дефер так же увеличивает х, но он отработает после return,
// в который скопировалась переменная х на момент обьявления в коде (т.е. х = 1)
func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
