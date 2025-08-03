.PHONY: build clean dev

build:
	npm install
	npm run build
	@echo "✅ Frontend built to: ./build/"

clean:
	rm -rf build/ node_modules/.vite/ .svelte-kit/

dev:
	npm run dev -- --host --port 5173