package main

import (
	"fmt"
)

func detectType(x interface{}) string {
	switch v := x.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	case chan int:
		return "chan int"
	case chan string:
		return "chan string"
	case chan bool:
		return "chan bool"
	case chan interface{}:
		return "chan interface{}"
	default:
		return fmt.Sprintf("unknown type %T", v)
	}
}

func main() {
	fmt.Println(detectType(42))
	fmt.Println(detectType("hello"))
	fmt.Println(detectType(true))
	fmt.Println(detectType(make(chan int)))
	fmt.Println(detectType(make(chan string)))
	fmt.Println(detectType(make(chan bool)))
	fmt.Println(detectType(3.14)) // тут краш ибо float64 не определялм
}
