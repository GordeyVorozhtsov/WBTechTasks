package main

import "fmt"

// интерфейс, с которым работает клиент
type Animal interface {
	MakeSound()
}

type Dog struct{}

func (d *Dog) Bark() {
	fmt.Println("Гав-гав!")
}

type DogAdapter struct {
	dog *Dog
}

// делает возможным вывод Bark()
func (a *DogAdapter) MakeSound() {
	a.dog.Bark()
}

func main() {
	var animal Animal
	animal = &DogAdapter{dog: &Dog{}}
	animal.MakeSound() // Гав-гав!
}
