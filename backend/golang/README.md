# Quotes API - Go Backend

Simple Go API for searching quotes using ParadeDB full-text search.

## Setup

1. Create `.env` symlink from base folder:
   ```bash
   ln -sf ../../.env .env
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

## Usage

- `make dev` - Start server with live reload (recommended for development)
- `make dev > server.log 2>&1 &` - Start server with live reload, log to file (background)
- `make run` - Start the server on port 8080
- `make test` - Run integration tests (server must be running)

## Endpoints

- `GET /health` - Health check
- `GET /api/search?q=life` - Search quotes

## Tech Stack

- Go standard library (`net/http`)
- `pgx/v5` - PostgreSQL driver
- `rs/cors` - CORS middleware
- ParadeDB BM25 search indexes