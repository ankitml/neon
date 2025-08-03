.PHONY: help build build-frontend build-backend build-docker clean deploy

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build everything (frontend + backend + docker)"
	@echo "  build-frontend - Build SvelteKit static files"
	@echo "  build-backend  - Build Go binary"
	@echo "  build-docker   - Build and tag Docker image"
	@echo "  clean          - Clean all build artifacts"
	@echo "  deploy         - Build and push to registry"

# Build everything
build: build-frontend build-backend build-docker

# Call subfolder makefiles
build-frontend:
	@echo "Building frontend..."
	$(MAKE) -C frontend build

build-backend:
	@echo "Building backend..."
	$(MAKE) -C backend/golang build

# Build Docker image using artifacts
build-docker: 
	@echo "Building Docker image..."
	docker build -t ghcr.io/ankitml/neon/quotes-app:latest .
	@echo "✅ Docker image: ghcr.io/ankitml/neon/quotes-app:latest"

# Clean all artifacts
clean:
	$(MAKE) -C frontend clean
	$(MAKE) -C backend/golang clean
	docker rmi ghcr.io/ankitml/neon/quotes-app:latest 2>/dev/null || true

# Build and push to registry
deploy: build
	docker push ghcr.io/ankitml/neon/quotes-app:latest
	@echo "✅ Pushed to registry: ghcr.io/ankitml/neon/quotes-app:latest"