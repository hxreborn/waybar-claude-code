.PHONY: build install clean fmt lint test run
.DEFAULT_GOAL := build

GO ?= go
BINARY ?= waybar-claude-code
CMD_DIR := ./cmd/waybar-claude-code
INSTALL_DIR ?= $(HOME)/.config/waybar/modules
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w -X main.version=$(VERSION)
RUN_INTERVAL ?= 0
STATICCHECK := $(shell command -v staticcheck)
GOTESTFLAGS ?= -race -cover

build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD_DIR)

install: build
	install -Dm755 $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY)"
	@echo "Restart Waybar to load the module: pkill -SIGUSR2 waybar"

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...
	@if [ -n "$(STATICCHECK)" ]; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed; skipping"; \
	fi

test:
	$(GO) test $(GOTESTFLAGS) ./...

run:
	CLAUDE_INTERVAL_SEC=$(RUN_INTERVAL) $(GO) run $(CMD_DIR)

clean:
	rm -f $(BINARY)
