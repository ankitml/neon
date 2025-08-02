.PHONY: help run dev test clean deps

# Default target
help:
	@echo "Available commands:"
	@echo "  run    - Start the Go server"
	@echo "  dev    - Start the Go server with live reload (recommended)"
	@echo "  test   - Run integration tests (server must be running)"
	@echo "  deps   - Install/update Go dependencies"
	@echo "  clean  - Clean Go module cache"

# Load environment and run server
run:
	@if [ ! -f .env ]; then echo "Error: .env file not found"; exit 1; fi
	@echo "Starting quotes API server..."
	go run main.go handlers.go

# Development with live reload
dev:
	@if [ ! -f .env ]; then echo "Error: .env file not found"; exit 1; fi
	@echo "Starting quotes API server with live reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# Run integration tests
test:
	@echo "Running integration tests..."
	@echo "Note: Make sure the server is running on localhost:8080"
	go test ./test -v

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean module cache
clean:
	go clean -modcache
	go mod tidy