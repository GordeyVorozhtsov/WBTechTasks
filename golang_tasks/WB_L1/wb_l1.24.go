package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

func newPoint(x, y float64) Point {
	return Point{x, y}
}

func (p Point) distance(p2 Point) float64 {
	dx := p.x - p2.x
	dy := p.y - p2.y
	return math.Sqrt(dx*dx + dy*dy) //формула нахождения расстояния между точками корень суммы квадратов (х1-х2) и (у1-у2)
}
func main() {
	p1 := newPoint(4, 5)
	p2 := newPoint(2, 3)
	d := p1.distance(p2)
	fmt.Println(d)
}
