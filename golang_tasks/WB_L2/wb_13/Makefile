TEST_TARGETS := $(shell grep -E '^test_[a-zA-Z0-9_-]+:' $(MAKEFILE_LIST) | cut -d ':' -f 1)

TEST_FILE := fixture-assets/data.txt

.PHONY: all build $(TEST_TARGETS) test-all

all: build

build:
	go build 

test-all: $(TEST_TARGETS)

test_field_number:
	./cut -f 1-3,5 $(TEST_FILE)

test_separator:
	./cut -f 1 -s $(TEST_FILE)
	
test_delimetr:
	./cut -d " " -f 2,4 $(TEST_FILE)