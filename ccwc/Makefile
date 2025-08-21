BINARY_NAME=ccwc
VERSION ?= $(shell git describe --tags --always --dirty)

LDFLAGS=-ldflags "-X 'github.com/suhalvemu/ccwc/cmd.version=$(VERSION)' -X 'main.version=$(VERSION)'"

.PHONY: all build test clean run

all: test build

test:
	@echo "==> Running tests..."
	go test ./...

build: clean
	@echo "==> Building binary..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -rf bin/
