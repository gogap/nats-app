.PHONY: build run clean test deps fmt lint

# Application name
APP_NAME := nats-client
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)

# Default target
all: build

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run tests
test:
	go test -v ./...

# Build the application
build: deps fmt
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .

# Build for multiple platforms
build-all: deps fmt
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME)-windows-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME)-darwin-arm64 .

# Run the application
run: deps
	go run -ldflags "$(LDFLAGS)" .

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -f $(APP_NAME)-*
	go clean

# Development mode with hot reload
dev:
	air

# Show help
help:
	@echo "Available targets:"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  test        - Run tests"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  run         - Run the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  dev         - Run in development mode"
	@echo "  help        - Show this help" 