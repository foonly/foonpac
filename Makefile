.PHONY: build dev prepare install clean test format

# Go build flags
LDFLAGS := -s -w -X github.com/foonly/foonpac/internal/config.AppVersion=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build: clean prepare bin/foonpac

bin/foonpac:
	mkdir -p bin
	go build -ldflags="$(LDFLAGS)" -o bin/foonpac .

dev: clean prepare
	mkdir -p bin
	go build -o bin/foonpac .

prepare:
	go mod download
	go mod tidy

install: build
	mkdir -p ~/.local/bin
	install -m 755 bin/foonpac ~/.local/bin/foonpac
	# Bash completion
	mkdir -p ~/.local/share/bash-completion/completions
	bin/foonpac completion bash > ~/.local/share/bash-completion/completions/foonpac
	# Zsh completion
	mkdir -p ~/.local/share/zsh/site-functions
	bin/foonpac completion zsh > ~/.local/share/zsh/site-functions/_foonpac
	# Fish completion
	mkdir -p ~/.local/share/fish/vendor_completions.d
	bin/foonpac completion fish > ~/.local/share/fish/vendor_completions.d/foonpac.fish

clean:
	rm -rf bin
	go clean

test:
	go test ./...

format:
	go fmt ./...
