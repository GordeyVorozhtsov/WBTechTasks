package main

import (
	"fmt"
	"math"
)

func main() {
	temperatureList := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	m := make(map[int][]float64)
	fmt.Println(splitingTemperature(m, temperatureList))
}
func splitingTemperature(m map[int][]float64, temperatureList []float64) map[int][]float64 {
	for _, e := range temperatureList {
		ceil := int(math.Ceil(e/10) * 10)
		if e > 0 {
			ceil = int(math.Floor(e/10) * 10)
		}
		m[ceil] = append(m[ceil], e)
	}
	return m
}
