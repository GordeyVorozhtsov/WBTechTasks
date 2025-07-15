package main

import "fmt"

func createHugeString(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = 'a'
	}
	return string(b)
}

var justString string

func someFunc() {
	v := createHugeString(1 << 10)       // здесь может быть любая большяа цифра и может много весить
	justString = string([]byte(v[:100])) // поэтмоу мы переопределяем justString и выдаем ей конкретно 100 элементов, а не делаем ее ссылкой на весь массив(не безопасно)
}

func main() {
	someFunc()
	fmt.Println(len(justString)) // 100
}
