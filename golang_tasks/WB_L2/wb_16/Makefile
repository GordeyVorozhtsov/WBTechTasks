TEST_TARGETS := $(shell grep -E '^test_[a-zA-Z0-9_-]+:' $(MAKEFILE_LIST) | cut -d ':' -f 1)

.PHONY: all build $(TEST_TARGETS) test-all

all: build

build:
	go build 

clean:
	rm -rf artifacts/*

test-all: $(TEST_TARGETS)

test-dzen:
	./wget dzen.ru

test-vk:
	./wget vk.com -t 1

test-neverssl:
	./wget neverssl.com -d 3 -j 10 -t 10

	
