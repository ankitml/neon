.PHONY: help install dev build start lint clean stop logs

# Default target
help:
	@echo "Available commands:"
	@echo "  install - Install dependencies using pnpm"
	@echo "  dev     - Start the Next.js development server in the background"
	@echo "  stop    - Stop the development server"
	@echo "  logs    - View the development server logs"
	@echo "  build   - Build the application for production"
	@echo "  start   - Start the production server"
	@echo "  lint    - Lint the codebase"
	@echo "  clean   - Remove node_modules and .next directories"

# Install dependencies
install:
	@echo "Installing dependencies..."
	pnpm install

# Development server
dev:
	@command -v pnpm >/dev/null 2>&1 || { echo >&2 "pnpm is not installed. Please install it first."; exit 1; }
	@touch dev-server.log
	@echo "Starting development server in the background..."
	pnpm dev > dev-server.log 2>&1 &

# Stop development server
stop:
	@echo "Stopping development server..."
	@pkill -f "next dev"

# View logs
logs:
	@echo "Tailing development server logs..."
	@tail -f dev-server.log

# Build for production
build:
	@echo "Building for production..."
	pnpm build

# Production server
start:
	@echo "Starting production server..."
	pnpm start

# Lint codebase
lint:
	@echo "Linting codebase..."
	pnpm lint

# Clean project
clean:
	@echo "Cleaning project..."
	rm -rf node_modules .next dev-server.log
