# Makefile for building and obfuscating go implant

BIN_NAME?=implant

all: setup deps build

setup:
	go install mvdan.cc/garble@latest
	
deps:
	go mod tidy

build:
	$(shell go env GOPATH)/bin/garble -literals -tiny build -ldflags="-s -w" -o $(BIN_NAME)
	