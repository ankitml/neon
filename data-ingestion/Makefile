.PHONY: help setup install clean migrate upgrade downgrade revision current history test lint format ingest

# Default target
help:
	@echo "Available commands:"
	@echo "  setup      - Create virtual environment"
	@echo "  install    - Install/update dependencies"
	@echo "  clean      - Remove virtual environment and cache files"
	@echo "  migrate    - Run all pending migrations"
	@echo "  upgrade    - Alias for migrate"
	@echo "  downgrade  - Downgrade database by one revision"
	@echo "  revision   - Create a new migration (use MSG='description')"
	@echo "  current    - Show current migration"
	@echo "  history    - Show migration history"
	@echo "  test       - Run tests (if any)"
	@echo "  lint       - Run code linting"
	@echo "  format     - Format code"
	@echo "  ingest     - Run quotes ingestor (use FILE=path/to/quotes.json)"

# Virtual environment setup
.venv/pyvenv.cfg:
	python3 -m venv .venv
	@echo "Virtual environment created in .venv"

setup: .venv/pyvenv.cfg

# Install dependencies using system pip3 with --break-system-packages
install:
	pip3 install -r requirements.txt --break-system-packages

# Migration commands
migrate:
	alembic upgrade head

upgrade: migrate

downgrade:
	alembic downgrade -1

revision:
	@if [ -z "$(MSG)" ]; then \
		echo "Error: Please provide a message with MSG='description'"; \
		echo "Example: make revision MSG='add user table'"; \
		exit 1; \
	fi
	alembic revision -m "$(MSG)"

current:
	alembic current

history:
	alembic history

# Development commands
test: .venv/lib/python*/site-packages/alembic
	@echo "No tests configured yet"

lint: .venv/lib/python*/site-packages/alembic
	@echo "No linting configured yet"

format: .venv/lib/python*/site-packages/alembic
	@echo "No formatting configured yet"

# Data ingestion
ingest:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Please provide a JSON file with FILE='path/to/quotes.json'"; \
		echo "Example: make ingest FILE='../quotes.json'"; \
		exit 1; \
	fi
	python3 -m src.data_ingestion.ingestor $(FILE)

# Cleanup
clean:
	rm -rf .venv
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
	@echo "Cleaned up virtual environment and cache files"