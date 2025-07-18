package main

import (
	"fmt"
	"os"
)

// классическая проверка на то что если значение может быть nil но значение переменной ссылочный тип и не является nil
func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
