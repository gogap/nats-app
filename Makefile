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

# Build development version (with embedded fonts)
build: fmt
	go build -ldflags="-X main.Version=dev -X main.BuildTime=$(shell date -u +'%Y-%m-%d_%H:%M:%S_UTC') -X main.GoVersion=$(shell go version | awk '{print $$3}')" -o $(APP_NAME) .

# Build with specific version
build-release: fmt
	go build -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u +'%Y-%m-%d_%H:%M:%S_UTC') -X main.GoVersion=$(shell go version | awk '{print $$3}')" -o $(APP_NAME) .

# Development build with live-reload (requires air)
build-dev:
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Run the application
run: deps
	go run -ldflags "$(LDFLAGS)" .

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -f nats-client
	rm -rf dist/
	rm -rf fyne-cross/

# Show help
help:
	@echo "Available targets:"
	@echo "  deps        - Install Go dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  test        - Run tests"
	@echo "  build       - Build the application (with embedded fonts)"
	@echo "  build-release - Build with release version"
	@echo "  build-dev   - Build with live-reload"
	@echo "  run         - Run the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  dev         - Run in development mode"
	@echo "  help        - Show this help" 