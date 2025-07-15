package main

import (
	"fmt"
)

func reverseString(x string) string {
	res := []rune(x)
	for i := 0; i < len(res)/2; i++ {
		res[i], res[len(res)-1-i] = res[len(res)-1-i], res[i]
	}
	return string(res)
}

func main() {
	fmt.Println(reverseString("главрыба"))
}
