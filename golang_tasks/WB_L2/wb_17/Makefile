TEST_TARGETS := $(shell grep -E '^test_[a-zA-Z0-9_-]+:' $(MAKEFILE_LIST) | cut -d ':' -f 1)

.PHONY: all build $(TEST_TARGETS) test-all

all: build

build:
	go build 

clean:
	rm -rf artifacts/*

test-all: $(TEST_TARGETS)

test_echo_socket:
	./telnet tcpbin.com -p 4242

test_google_socket:
	./telnet google.com -p 80  -t 10