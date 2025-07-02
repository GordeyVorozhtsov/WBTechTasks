package main

import (
	"fmt"
)

func setBit(num int64, i uint, bit uint) int64 {
	if bit == 1 {
		// установить i-й бит в 1
		return num | (1 << i)
	}
	// установить i-й бит в 0
	return num &^ (1 << i)
}

func main() {
	var num int64 = 5 // 0101 в двоичном виде
	fmt.Printf("Исходное число: %d (%.4b)\n", num, num)

	// установить 1-й бит (нумерация с 0) в 0
	num = setBit(num, 1, 0)
	fmt.Printf("После установки 1-го бита в 0: %d (%.4b)\n", num, num)

	// установать 2-й бит в 1
	num = setBit(num, 2, 1)
	fmt.Printf("После установки 2-го бита в 1: %d (%.4b)\n", num, num)
}
