.PHONY: build install clean fmt lint test run

BINARY ?= waybar-claude-code
CMD_DIR := ./cmd/waybar-claude-code
INSTALL_DIR ?= $(HOME)/.config/waybar/modules
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
STATICCHECK := $(shell command -v staticcheck)
RUN_INTERVAL ?= 0

build:
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "-s -w -X main.version=$(VERSION)" \
		-o $(BINARY) \
		$(CMD_DIR)

install: build
	install -Dm755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY)"
	@echo "Restart Waybar to load the module: pkill -SIGUSR2 waybar"

fmt:
	go fmt ./...

lint:
	go vet ./...
	@if [ -n "$(STATICCHECK)" ]; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed; skipping"; \
	fi

test:
	go test -race -cover ./...

run:
	CLAUDE_INTERVAL_SEC=$(RUN_INTERVAL) go run $(CMD_DIR)

clean:
	rm -f $(BINARY)
