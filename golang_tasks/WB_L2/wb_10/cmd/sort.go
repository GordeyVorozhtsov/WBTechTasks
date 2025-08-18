package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	keyColumn            int  // Номер колонки, по которой сортируем (-k)
	numericSort          bool // сортировать как числа (-n)
	reverseSort          bool // сортировать в обратном порядке (-r)
	unique               bool // выводить только уникальные строки (-u)
	monthSort            bool // сортировать месяцы по порядку (-M)
	ignoreTrailingBlanks bool // игнорировать хвостовые пробелы (-b)
	checkOrder           bool // проверить отсортирован ли вход (-c)
	humanNumeric         bool // сортировать по размеру с суффиксами (-H)
)

var sortCmd = &cobra.Command{
	Use:   "sort [file][flags]",
	Short: "Sort lines of text with various options",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSort,
}

func init() {
	// Инициализация флагов для команды sort
	sortCmd.Flags().IntVarP(&keyColumn, "key", "k", 0, "Sort by column number (starting from 1). Default: entire line")
	sortCmd.Flags().BoolVarP(&numericSort, "numeric", "n", false, "Sort by numeric value")
	sortCmd.Flags().BoolVarP(&reverseSort, "reverse", "r", false, "Sort in reverse order")
	sortCmd.Flags().BoolVarP(&unique, "unique", "u", false, "Output only unique lines")
	sortCmd.Flags().BoolVarP(&monthSort, "month", "M", false, "Sort by month name (Jan, Feb, ...)")
	sortCmd.Flags().BoolVarP(&ignoreTrailingBlanks, "ignore-trailing-blanks", "b", false, "Ignore trailing blanks")
	sortCmd.Flags().BoolVarP(&checkOrder, "check", "c", false, "Check whether input is sorted; print message if not")
	sortCmd.Flags().BoolVarP(&humanNumeric, "human-numeric", "H", false, "Sort by human readable sizes (e.g. 1K, 2M)")
	rootCmd.AddCommand(sortCmd) // Добавляем команду в корень
}

func runSort(cmd *cobra.Command, args []string) error {
	lines, err := readInputLines(args)
	if err != nil {
		return err
	}

	// обрезаем у каждой строки справа пробелы и табы
	if ignoreTrailingBlanks {
		for i := range lines {
			lines[i] = strings.TrimRight(lines[i], " \t")
		}
	}

	// Функция для вытаскивания ключа сортировки из строки
	getKey := func(line string) string {
		if keyColumn <= 0 {
			// Если колонки не выбрали, сортируем по всей строке
			return line
		}
		// Разбиваем строку на слова
		fields := strings.Fields(line)
		if keyColumn > len(fields) { // Если нет такой колонки в строке
			return ""
		}
		// Берём нужное поле
		return fields[keyColumn-1]
	}

	// Функция для сравнения двух ключей в зависимости от типа сортировки
	compare := func(a, b string) int {
		switch {
		case monthSort:
			return compareMonths(a, b) // Сравниваем как месяцы
		case numericSort:
			return compareNumbers(a, b) // Сравниваем как числа
		case humanNumeric:
			return compareHumanSizes(a, b) // Сравниваем как "человеческие" размеры с K, M, G
		default:
			return strings.Compare(a, b) // Обычное строковое сравнение
		}
	}

	// Функция, определяющая порядок для сортировки (меньше ли i-й элемент, чем j-й)
	less := func(i, j int) bool {
		a := getKey(lines[i])
		b := getKey(lines[j])
		cmp := compare(a, b)
		if reverseSort {
			return cmp > 0 // Если флаг обратной сортировки меняем логику
		}
		return cmp < 0
	}

	// Если включена проверка отсортированности
	if checkOrder {
		for i := 1; i < len(lines); i++ {
			if less(i, i-1) {
				// Если нашли нарушение порядка сообщаем и прерываем сортировку
				fmt.Fprintln(os.Stderr, "Input is not sorted")
				return fmt.Errorf("input is not sorted")
			}
		}
		return nil
	}

	// Собственно, сортируем
	sort.SliceStable(lines, less)

	// Если надо используем unique sort
	if unique {
		lines = uniqueLines(lines)
	}

	// Выводим результаты построчно
	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

// Читаем строки либо из файла, либо из стандартного ввода
func readInputLines(args []string) ([]string, error) {
	var reader io.Reader
	if len(args) == 0 {
		reader = os.Stdin
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader = f
	}

	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Убираем подряд идущие одинаковые строки, оставляем только 1-й экземпляр
func uniqueLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	result := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}
	return result
}

// Сопоставляем названию месяца число (для сортировки по месяцам)
var monthMap = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
	"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
	"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

// Функция сравнения двух месяцев по номеру
func compareMonths(a, b string) int {
	ma, oka := monthMap[fixMonthLetter(a)]
	mb, okb := monthMap[fixMonthLetter(b)]
	if !oka && !okb {
		// Если ничё не поняли, просто строкой сравним
		return strings.Compare(a, b)
	}
	if !oka {
		return -1 // a непонятен считаем меньше
	}
	if !okb {
		return 1 // b непонятен считаем больше
	}
	if ma < mb {
		return -1
	} else if ma > mb {
		return 1
	}
	return 0
}

// Нормализует строку месяца делает 1 букву заглавной, остальные строчными
func fixMonthLetter(s string) string {
	if len(s) < 1 {
		return s
	}
	if len(s) < 3 {
		s = strings.ToLower(s)
		return strings.ToUpper(s[:1]) + s[1:]
	}
	s = strings.ToLower(s[:3])
	return strings.ToUpper(s[:1]) + s[1:]
}

// Сравниваем две строки как числа
func compareNumbers(a, b string) int {
	fa, errA := strconv.ParseFloat(a, 64)
	fb, errB := strconv.ParseFloat(b, 64)
	if errA != nil && errB != nil {
		return strings.Compare(a, b)
	}
	if errA != nil {
		return -1 // a не число, значит меньше
	}
	if errB != nil {
		return 1 // b не число, значит меньше
	}
	switch {
	case fa < fb:
		return -1
	case fa > fb:
		return 1
	default:
		return 0
	}
}

// Сравнение "человеческих" размеров с суффиксами K, M, G, T
func compareHumanSizes(a, b string) int {
	fa, errA := parseHumanSize(a)
	fb, errB := parseHumanSize(b)

	if errA != nil && errB != nil {
		return strings.Compare(a, b)
	}
	if errA != nil {
		return -1
	}
	if errB != nil {
		return 1
	}
	switch {
	case fa < fb:
		return -1
	case fa > fb:
		return 1
	default:
		return 0
	}
}

// Парсим строку с размером типа "1K", "5M", "10G", "1000" в float64 байтовый эквивалент
func parseHumanSize(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size string")
	}
	mult := 1.0
	last := s[len(s)-1]
	numberPart := s

	switch last {
	case 'K', 'k':
		mult = 1024
		numberPart = s[:len(s)-1]
	case 'M', 'm':
		mult = 1024 * 1024
		numberPart = s[:len(s)-1]
	case 'G', 'g':
		mult = 1024 * 1024 * 1024
		numberPart = s[:len(s)-1]
	case 'T', 't':
		mult = 1024 * 1024 * 1024 * 1024
		numberPart = s[:len(s)-1]
	}

	f, err := strconv.ParseFloat(numberPart, 64)
	if err != nil {
		return 0, err
	}
	return f * mult, nil
}
