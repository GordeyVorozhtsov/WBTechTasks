package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	fieldsStr string
	delimiter string
	separated bool
)

func init() {
	// Добавляем флаги к rootCmd
	rootCmd.Flags().StringVarP(&fieldsStr, "fields", "f", "", "Select only these fields (e.g. 1,3-5) (required)")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "\t", "Use a different delimiter character")
	rootCmd.Flags().BoolVarP(&separated, "separated", "s", false, "Only output lines containing delimiter")

	// Делаем флаг fields обязательным
	rootCmd.MarkFlagRequired("fields")

	// Настраиваем функцию Run
	rootCmd.Run = runCut
}

func runCut(cmd *cobra.Command, args []string) {
	// Парсим список полей
	fieldIndices, err := parseFields(fieldsStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing fields: %v\n", err)
		os.Exit(1)
	}

	// Определяем источник ввода
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	// Читаем строки из файла или stdin
	lines, err := readInputLines(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Обрабатываем каждую строку
	for _, line := range lines {
		processLine(line, fieldIndices)
	}
}

func parseFields(fieldsStr string) ([]int, error) {
	var fields []int
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			// Обрабатываем диапазон (например, "3-5")
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil || start < 1 {
				return nil, fmt.Errorf("invalid start field: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil || end < 1 {
				return nil, fmt.Errorf("invalid end field: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("start field cannot be greater than end field: %s", part)
			}

			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			// Обрабатываем отдельное поле
			field, err := strconv.Atoi(part)
			if err != nil || field < 1 {
				return nil, fmt.Errorf("invalid field: %s", part)
			}
			fields = append(fields, field)
		}
	}

	return fields, nil
}

func processLine(line string, fieldIndices []int) {
	// Если включен флаг -s и строка не содержит разделитель, пропускаем её
	if separated && !strings.Contains(line, delimiter) {
		return
	}

	// Разбиваем строку на поля
	fields := strings.Split(line, delimiter)

	// Собираем только нужные поля
	var outputFields []string
	for _, idx := range fieldIndices {
		// Индексы начинаются с 1, а в массиве - с 0
		if idx-1 < len(fields) {
			outputFields = append(outputFields, fields[idx-1])
		}
	}

	// Выводим результат
	if len(outputFields) > 0 {
		fmt.Println(strings.Join(outputFields, delimiter))
	}
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
