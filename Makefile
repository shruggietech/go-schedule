# go-schedule, developer tasks
# The Wails desktop uses its separate module, Node frontend and native WebView.
# The daemon and CLI are cgo-free. Canonical verification covers both modules.

GO        ?= go
SH        ?= sh
LDFLAGS   ?=
WAILS_VERSION = v2.17.0

.PHONY: all verify fmt fmt-check vet lint test test-race test-gui cover docs-check automation-check bench build build-daemon build-cli build-gui tidy clean

all: fmt vet test build

# Mutating convenience target. Use fmt-check or verify for validation.
fmt:
	$(GO) fmt ./...

verify:
	$(SH) scripts/verify.sh all

fmt-check:
	$(SH) scripts/verify.sh format

vet:
	$(SH) scripts/verify.sh vet

lint:
	$(SH) scripts/verify.sh lint

test:
	$(GO) test -timeout=60m ./...

test-race:
	$(SH) scripts/verify.sh race

cover:
	$(SH) scripts/verify.sh coverage

bench:
	$(GO) test -timeout=60m -bench=. -benchmem ./internal/engine/...

build: build-daemon build-cli

build-daemon:
	$(GO) build -o bin/goschedd ./cmd/goschedd

build-cli:
	$(GO) build -o bin/gosched ./cmd/gosched

# Wails needs Node and platform WebView development prerequisites.
build-gui:
	cd desktop && $(GO) run github.com/wailsapp/wails/v2/cmd/wails@$(WAILS_VERSION) build -clean -tags webkit2_41

# Wails uses the windowless Windows linker mode for its production desktop.
build-gui-windows:
	cd desktop && $(GO) run github.com/wailsapp/wails/v2/cmd/wails@$(WAILS_VERSION) build -clean -platform windows/amd64

test-gui:
	$(SH) scripts/verify.sh gui

docs-check:
	$(SH) scripts/verify.sh docs

automation-check:
	$(SH) scripts/verify.sh automation

tidy:
	$(GO) mod tidy

clean:
	rm -rf bin coverage.out
