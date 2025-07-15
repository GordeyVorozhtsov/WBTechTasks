package main

import (
	"fmt"
)

func discinctElement(a []string) []string {
	res := []string{}
	m := make(map[string]int)
	for _, elem := range a {
		m[elem]++
	}
	for key := range m {
		res = append(res, key) // не считаю что это правльно по условию но это подходит к примеру результата задачи
	}
	return res
}
func main() {
	fmt.Println(discinctElement([]string{"cat", "cat", "dog", "cat", "tree"}))
}

// это реализация конкретно уникальных элементов из массива, а пример в условии не корректен
// func discinctElement(a []string) []string {
// 	res := []string{}
// 	m := make(map[string]int)
// 	for _, elem := range a {
// 		m[elem]++
// 	}
// 	for key := range m {
// 		if m[key] == 1 {
// 			res = append(res, key)
// 		}
// 	}
// 	return res
// }
