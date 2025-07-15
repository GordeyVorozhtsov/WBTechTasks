package main

import (
	"fmt"
	"strings"
)

func uniqueElems(arr string) bool {
	strings.ToLower(arr)
	uq := make(map[rune]int) // мапа с уникальными элементами(количество == 1)
	for _, elem := range arr {
		uq[elem]++
		if uq[elem] > 1 {
			return false
		}
	}
	return true
}
func main() {
	fmt.Println(uniqueElems("abcd"))
	fmt.Println(uniqueElems("abCdefAf"))
	fmt.Println(uniqueElems("aabcd"))

}
