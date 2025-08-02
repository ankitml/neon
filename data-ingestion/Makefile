.PHONY: help setup install clean migrate upgrade downgrade revision current history test lint format

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

# Virtual environment setup
.venv/pyvenv.cfg:
	python3 -m venv .venv
	@echo "Virtual environment created in .venv"

setup: .venv/pyvenv.cfg

# Install dependencies (updates when requirements.txt changes)
.venv/lib/python*/site-packages/alembic: .venv/pyvenv.cfg requirements.txt
	.venv/bin/pip install --upgrade pip
	.venv/bin/pip install -r requirements.txt
	@touch .venv/lib/python*/site-packages/alembic

install: .venv/lib/python*/site-packages/alembic

# Migration commands
migrate: .venv/lib/python*/site-packages/alembic
	.venv/bin/alembic upgrade head

upgrade: migrate

downgrade: .venv/lib/python*/site-packages/alembic
	.venv/bin/alembic downgrade -1

revision: .venv/lib/python*/site-packages/alembic
	@if [ -z "$(MSG)" ]; then \
		echo "Error: Please provide a message with MSG='description'"; \
		echo "Example: make revision MSG='add user table'"; \
		exit 1; \
	fi
	.venv/bin/alembic revision --autogenerate -m "$(MSG)"

current: .venv/lib/python*/site-packages/alembic
	.venv/bin/alembic current

history: .venv/lib/python*/site-packages/alembic
	.venv/bin/alembic history

# Development commands
test: .venv/lib/python*/site-packages/alembic
	@echo "No tests configured yet"

lint: .venv/lib/python*/site-packages/alembic
	@echo "No linting configured yet"

format: .venv/lib/python*/site-packages/alembic
	@echo "No formatting configured yet"

# Cleanup
clean:
	rm -rf .venv
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
	@echo "Cleaned up virtual environment and cache files"