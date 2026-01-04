.PHONY: build build-daemon build-tui clean help

help:
	@echo "BanForge build targets:"
	@echo "  make build         - Build both daemon and TUI"
	@echo "  make build-daemon  - Build only daemon"
	@echo "  make build-tui     - Build only TUI"
	@echo "  make clean         - Remove binaries"
	@echo "  make test          - Run tests"	

build: build-daemon build-tui
	@echo "✅ Build complete!"

build-daemon:
	@mkdir -p bin
	go mod tidy
	go build -o bin/banforge ./cmd/banforge

build-tui:
	@mkdir -p bin
	go build -o bin/banforge-tui ./cmd/banforge-tui

clean:
	rm -rf bin/

test:
	go test ./...
