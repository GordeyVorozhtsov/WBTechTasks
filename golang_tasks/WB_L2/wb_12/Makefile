# Все цели
TEST_TARGETS := $(shell grep -E '^test_[a-zA-Z0-9_-]+:' $(MAKEFILE_LIST) | cut -d ':' -f 1)

# Путь к тестовому файлу
TEST_FILE := fixture-assets/data.txt

.PHONY: all build $(TEST_TARGETS) test-all

all: build

build:
	go build 

test-all: $(TEST_TARGETS)

test_grep_fixed_string:
	./grep grep -F "ERROR" $(TEST_FILE)

test_grep_ignore_case:
	./grep grep -i "test" $(TEST_FILE)

test_grep_count:
	./grep grep -c "the" $(TEST_FILE)

test_grep_invert_match:
	./grep grep -v "the" $(TEST_FILE)

test_grep_context_2:
	./grep grep -C 2 "ERROR" $(TEST_FILE)

test_grep_after_3:
	./grep grep -A 3 "ERROR" $(TEST_FILE)

test_grep_before_1:
	./grep grep -B 1 "ERROR" $(TEST_FILE)

test_grep_line_number:
	./grep grep -n "grep" $(TEST_FILE)

test_grep_multiple_flags:
	./grep grep -n -i "test" $(TEST_FILE)

test_grep_fixed_with_ignore_case:
	./grep grep -F -i "error" $(TEST_FILE)