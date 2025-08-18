package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	noLetter = errors.New("no letters in str")
)

func UnpackingString(arr string) (string, error) {
	var res []string
	integer := "123456789"
	letters := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

	splitString := []rune(arr)
	escaped := false
	var lastChar rune = 0

	for idx := 0; idx < len(splitString); idx++ {
		elem := splitString[idx]

		if escaped {
			// Экранированный символ — добавляем как есть и запоминаем
			res = append(res, string(elem))
			lastChar = elem
			escaped = false
			continue
		}

		if elem == '\\' {
			escaped = true
			continue
		}

		// Если текущий символ цифра, то пытаемся повторить lastChar
		if strings.Contains(integer, string(elem)) {
			if lastChar == 0 {
				// Нет предыдущего символа для повторения -> ошибка
				return "", noLetter
			}
			num, _ := strconv.Atoi(string(elem))
			// Мы уже добавили 1 раз lastChar, нужно добавить num раз
			for i := 1; i < num; i++ {
				res = append(res, string(lastChar))
			}
		}

		// Если буква просто добавляем и запоминаем
		if strings.Contains(letters, string(elem)) {
			res = append(res, string(elem))
			lastChar = elem
		}
	}

	resStr := strings.Join(res, "")
	if len(resStr) == 0 {
		return resStr, noLetter
	}
	return resStr, nil
}

func main() {
	fmt.Println(UnpackingString("a4bc2d5e")) // aaaabccddddde
	fmt.Println(UnpackingString("abcd"))     // abcd
	fmt.Println(UnpackingString("45"))       // error
	fmt.Println(UnpackingString(""))         // error
	fmt.Println(UnpackingString(`qwe\4\5`))  // qwe45
	fmt.Println(UnpackingString(`qwe\45`))   // qwe44444
}
