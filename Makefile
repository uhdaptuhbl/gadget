SHELL=/bin/bash
GOCMD=go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test

all: mods format lint

mods:
	$(GOCMD) mod tidy
	$(GOCMD) mod vendor
	$(GOCMD) mod verify

lint:
	golangci-lint run --timeout=5m --skip-dirs vendor ./...

format:
	find . -name \*.go -not -path vendor -not -path target -exec goimports -w {} \;

test:
	$(GOTEST) -v ./...
