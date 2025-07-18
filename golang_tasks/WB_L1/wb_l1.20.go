package main

import (
	"fmt"
	"strings"
)

func reverseString(arr string) string {
	res := strings.Split(arr, " ")

	for i := 0; i < len(res)/2; i++ {
		res[i], res[len(res)-1-i] = res[len(res)-1-i], res[i]
	}
	arr = strings.Join(res, " ")
	return arr
}
func main() {
	fmt.Println(reverseString("snow dog sun"))
}
