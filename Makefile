TOOLS_BIN=$(shell pwd)/.tools

.PHONY: install
install:
	@go mod download && \
		GO111MODULE=on GODEBUG=netdns=9 CGO_ENABLED=1  go build -ldflags="-s -w" -a -o $(TOOLS_BIN)/devtool

.PHONY: install_local
install_local:
	@go mod download && \
		GO111MODULE=on GODEBUG=netdns=9 CGO_ENABLED=1  go build -ldflags="-s -w" -a -o ~/.local/bin/devtool

.PHONY: setup
setup:
	@$(TOOLS_BIN)/devtool setup-lib

.PHONY: lint
lint:
	@$(TOOLS_BIN)/devtool lint

.PHONY: build
build:
	@$(TOOLS_BIN)/devtool build --arch=amd64

.PHONY: tests
tests:
	@$(TOOLS_BIN)/devtool test

.PHONY: ci
ci: install setup lint build tests

