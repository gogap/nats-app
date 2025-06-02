.PHONY: build run clean test deps fmt lint fyne-deps

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

# Install fyne packaging tool
fyne-deps: deps
	go install fyne.io/tools/cmd/fyne@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run tests
test:
	go test -v ./...

# Build the application using fyne package
build: fyne-deps fmt
	fyne package --name $(APP_NAME)

# Build for multiple platforms using fyne package
build-all: fyne-deps fmt
	fyne package --os linux --name $(APP_NAME)-linux-amd64
	fyne package --os windows --name $(APP_NAME)-windows-amd64
	fyne package --os darwin --name $(APP_NAME)-darwin-amd64
	fyne package --os darwin --name $(APP_NAME)-darwin-arm64

# Build using go build (for development/testing)
build-dev: deps fmt
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME)-dev .

# Run the application
run: deps
	go run -ldflags "$(LDFLAGS)" .

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -f $(APP_NAME)-*
	rm -rf $(APP_NAME)-*.app
	go clean

# Development mode with hot reload
dev:
	air

# Show help
help:
	@echo "Available targets:"
	@echo "  deps        - Install Go dependencies"
	@echo "  fyne-deps   - Install fyne packaging tool"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  test        - Run tests"
	@echo "  build       - Build the application using fyne package"
	@echo "  build-all   - Build for multiple platforms using fyne package"
	@echo "  build-dev   - Build using go build (for development)"
	@echo "  run         - Run the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  dev         - Run in development mode"
	@echo "  help        - Show this help" 