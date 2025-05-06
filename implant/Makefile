BINARY  := myapp

BIN_DIR := bin

GO      := go

GARBLE  := garble

PACKAGES := $(shell $(GO) list ./... | grep -v /vendor/)

.PHONY: all install-garble deps build clean

all: install-garble deps build

install-garble:
	@which $(GARBLE) >/dev/null 2>&1 || \
		$(GO) install mvdan.cc/garble@latest

deps:
	$(GO) mod tidy

build:
	@mkdir -p $(BIN_DIR)
	$(GARBLE) -literals -tiny build -o $(BIN_DIR)/$(BINARY) ./...