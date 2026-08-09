.PHONY: build build-daemon clean help install-man check-man

help:
	@echo "BanForge build targets:"
	@echo "  make build         - Build daemon binary"
	@echo "  make build-daemon  - Build only daemon"
	@echo "  make clean         - Remove binaries"
	@echo "  make test          - Run tests"
	@echo "  make install-man   - Install manpages to system"
	@echo "  make check-man     - Validate manpage syntax"	

build: build-daemon
	@echo "✅ Build complete!"

build-daemon:
	@mkdir -p bin
	go mod tidy
	go build -o bin/banforge ./cmd/banforge

clean:
	rm -rf bin/

test:
	go test ./...

test-cover:
	go test -cover ./...

lint:
	golangci-lint run --fix

check-man:
	@echo "Checking manpage syntax..."
	@man -l docs/man/banforge.1 > /dev/null && echo "✅ banforge.1 OK"
	@man -l docs/man/banforge.5 > /dev/null && echo "✅ banforge.5 OK"

install-man:
	@echo "Installing manpages..."
	install -d $(DESTDIR)/usr/share/man/man1
	install -d $(DESTDIR)/usr/share/man/man5
	install -m 644 docs/man/banforge.1 $(DESTDIR)/usr/share/man/man1/banforge.1
	install -m 644 docs/man/banforge.5 $(DESTDIR)/usr/share/man/man5/banforge.5
	@echo "✅ Manpages installed!"
