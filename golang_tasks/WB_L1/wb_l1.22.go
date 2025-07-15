package main

import (
	"fmt"
	"math/big"
)

// аргумент *big.Int ссыльный потому что под капотом структура и,
// чтобы не аллоцировать новую памать, проще передать так
func multiply(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Mul(a, b)
	return res
}

func divide(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Div(a, b)
	return res
}

func add(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Add(a, b)
	return res
}

func subtract(a, b *big.Int) *big.Int {
	res := new(big.Int)
	res.Sub(a, b)
	return res
}

func main() {
	a := big.NewInt(1234567890123456789)
	b := big.NewInt(987654321098765432)

	fmt.Println(multiply(a, b))
	fmt.Println(divide(a, b))
	fmt.Println(add(a, b))
	fmt.Println(subtract(a, b))
}
