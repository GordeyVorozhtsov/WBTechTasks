# Все цели, начинающиеся с test_
TEST_TARGETS := $(shell grep -E '^test_[a-zA-Z0-9_-]+:' Makefile | cut -d ':' -f 1)

# Путь к директории с примерами
EXAMPLES := fixture-assets/data.txt

.PHONY: cli-sorter $(TEST_TARGETS) test-all

build:
	go build cli-grep 

test-all: $(TEST_TARGETS)

test_cli_sorter_column_sort:
	./cli-sorter sort -k 2 $(EXAMPLES)

test_cli_sorter_num_sort:
	./cli-sorter sort $(EXAMPLES) -k 3 -n

test_cli_sorter_reverse:
	./cli-sorter sort -r $(EXAMPLES)

test_cli_sorter_unique_string:
	./cli-sorter sort $(EXAMPLES) -k 2 -u

test_cli_sorter_month:
	./cli-sorter sort $(EXAMPLES) -k 4 -M

test_cli_sorter_trailing_blanks:
	./cli-sorter sort $(EXAMPLES) -k 2 -b

test_cli_sorter_check_sort:
	./cli-sorter sort $(EXAMPLES) -k 1 -c

test_cli_sorter_human_numeric:
	./cli-sorter sort $(EXAMPLES) -k 5 -H

