package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/spf13/cobra"
)

var (
	stringsAfter   int  // -A вывести N строк после найденной
	stringsBefore  int  // -B вывести N строк до найденной
	stringsAround  int  // -C вывести N строк контекста вокруг
	sameStrings    bool // -c только количество совпадений
	ignoreRegistre bool // -i игнорировать фильтр
	reverseCurent  bool // -v инвертировать совпадения
	fixedString    bool // -F воспринимать шаблон как фиксированную строку
	stringNum      bool // -n выводить номер строки
)

var grepCmd = &cobra.Command{
	Use:   "grep [pattern] [file]",
	Short: "Search for patterns in text",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  runGrep,
}

func init() {
	grepCmd.Flags().IntVarP(&stringsAfter, "after", "A", 0, "Print N strings after match")
	grepCmd.Flags().IntVarP(&stringsBefore, "before", "B", 0, "Print N strings before match")
	grepCmd.Flags().IntVarP(&stringsAround, "context", "C", 0, "Print N strings around match")
	grepCmd.Flags().BoolVarP(&sameStrings, "count", "c", false, "Only print count of matching lines")
	grepCmd.Flags().BoolVarP(&ignoreRegistre, "ignore-case", "i", false, "Ignore case distinctions")
	grepCmd.Flags().BoolVarP(&reverseCurent, "invert-match", "v", false, "Select non-matching lines")
	grepCmd.Flags().BoolVarP(&fixedString, "fixed-strings", "F", false, "Interpret pattern as fixed string")
	grepCmd.Flags().BoolVarP(&stringNum, "line-number", "n", false, "Print line number with output")
	rootCmd.AddCommand(grepCmd)
}

func runGrep(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	var filename string
	if len(args) > 1 {
		filename = args[1]
	}

	lines, err := readInputLines(filename)
	if err != nil {
		return err
	}

	// Обработка флага -C
	if stringsAround > 0 {
		stringsBefore = stringsAround
		stringsAfter = stringsAround
	}

	// Подготовка паттерна для поиска
	var re *regexp.Regexp
	if fixedString {
		pattern = regexp.QuoteMeta(pattern)
	}
	if ignoreRegistre {
		pattern = "(?i)" + pattern
	}
	re, err = regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}

	// Поиск совпадений
	var matches []int
	for i, line := range lines {
		matched := re.MatchString(line)
		if reverseCurent {
			matched = !matched
		}
		if matched {
			matches = append(matches, i)
		}
	}

	// Обработка флага -c
	if sameStrings {
		fmt.Println(len(matches))
		return nil
	}

	// Сбор строк для вывода с учетом контекста
	outputLines := make(map[int]bool)
	for _, match := range matches {
		// Добавляем саму найденную строку
		outputLines[match] = true

		// Добавляем строки до
		for i := 1; i <= stringsBefore; i++ {
			if match-i >= 0 {
				outputLines[match-i] = true
			}
		}

		// Добавляем строки после
		for i := 1; i <= stringsAfter; i++ {
			if match+i < len(lines) {
				outputLines[match+i] = true
			}
		}
	}

	// Преобразуем в отсортированный список индексов
	var sortedIndices []int
	for k := range outputLines {
		sortedIndices = append(sortedIndices, k)
	}
	sort.Ints(sortedIndices)

	// Вывод результатов
	for _, idx := range sortedIndices {
		line := lines[idx]
		if stringNum {
			line = fmt.Sprintf("%d:%s", idx+1, line)
		}
		fmt.Println(line)
	}

	return nil
}

func readInputLines(filename string) ([]string, error) {
	var reader io.Reader

	if filename == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
