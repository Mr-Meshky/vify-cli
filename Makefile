BINARY_NAME=vify
BUILD_DIR=bin
VERSION?=1.0.0
LDFLAGS=-ldflags "-s -w -X 'github.com/Mr-Meshky/vify-cli/cmd.version=$(VERSION)'"

.PHONY: all build test clean install cross-build

all: test build

build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR) dist

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

cross-build:
	@mkdir -p $(BUILD_DIR)
	# macOS (Apple Silicon & Intel)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# Linux (x86_64 & ARM64)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	# Windows (x86_64)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
