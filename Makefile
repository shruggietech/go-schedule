# go-schedule — developer tasks
# The GUI binary (gosched-gui) requires a C toolchain + OpenGL (Fyne); the daemon
# and CLI are cgo-free. Targets below operate on the whole module.

GO        ?= go
SH        ?= sh
LDFLAGS   ?=
GUI_LDFLAGS_WINDOWS = -H windowsgui

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
	$(GO) test ./...

test-race:
	$(SH) scripts/verify.sh race

cover:
	$(SH) scripts/verify.sh coverage

bench:
	$(GO) test -bench=. -benchmem ./internal/engine/...

build: build-daemon build-cli

build-daemon:
	$(GO) build -o bin/goschedd ./cmd/goschedd

build-cli:
	$(GO) build -o bin/gosched ./cmd/gosched

# GUI: requires cgo + a C toolchain and OpenGL/X11 dev libraries (Fyne).
# On Windows add $(GUI_LDFLAGS_WINDOWS) so no console window appears.
build-gui:
	CGO_ENABLED=1 $(GO) build -o bin/gosched-gui ./cmd/gosched-gui

build-gui-windows:
	CGO_ENABLED=1 GOOS=windows $(GO) build -ldflags "$(GUI_LDFLAGS_WINDOWS)" -o bin/gosched-gui.exe ./cmd/gosched-gui

# Headless GUI tests run without a display or OpenGL (Fyne test driver).
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
