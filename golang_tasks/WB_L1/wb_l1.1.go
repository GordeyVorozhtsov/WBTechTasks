package main

import (
	"fmt"
)

type Human struct {
	Name string
}

func (h *Human) Helloing() {
	fmt.Printf("i m %s\n", h.Name)
}

type Action struct {
	Human
}

func main() {
	a := Action{Human: Human{Name: "WB_worker"}}
	a.Helloing()
}
