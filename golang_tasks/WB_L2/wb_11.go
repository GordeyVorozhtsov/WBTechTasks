package main
import(
	"fmt"
	"sort"
)
func annagrama(arr []string) map[string][]string {
    certArr := arr
    res := make(map[string][]string)

    for i := 0; i < len(certArr); {
        current := certArr[i]

        for j := i + 1; j < len(certArr); j++ {
        	// Если слово аннаграма, мы добавляем как ключ current и как значение append слова аннаграмму
            if isAnnagram(current, certArr[j]) {
                res[current] = append(res[current], certArr[j])
                // Удаляем элемент с индексом j
                certArr = append(certArr[:j], certArr[j+1:]...)
                j-- // сдвигаем индекс после удаления
            }
        }
        i++
    }

    // Сортируем строки по ключу карт
    for key := range res {
        sort.Strings(res[key])
    }

    return res
}

func isAnnagram(first, second string) bool{
	f := []rune(first)
	s := []rune(second)
	elems := make(map[rune]int)
	for _,e := range f{
		elems[e]++
	}
	for _,e := range s{
		if elems[e] == 0{
			return false
		}
	}
	return true 

}
func main(){
	str := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(annagrama(str))
}