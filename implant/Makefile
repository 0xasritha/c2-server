BINARY  := myapp

TMP_DIR := tmp

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
	@mkdir -p $(TMP_DIR)
	$(GARBLE) -literals -tiny build -o $(TMP_DIR)/$(BINARY) ./...