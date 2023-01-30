TOOLS_BIN=$(shell pwd)/.tools

.PHONY: install
install:
	go mod download
	go build -v -a -o $(TOOLS_BIN)/devtool

.PHONY: setup
setup:
	$(TOOLS_BIN)/devtool setup-lib

.PHONY: lint
lint:
	$(TOOLS_BIN)/devtool lint

.PHONY: build
build:
	$(TOOLS_BIN)/devtool build --arch=amd64

.PHONY: tests
tests:
	$(TOOLS_BIN)/devtool test

.PHONY: pre-commite
pre-commite: setup lint build tests

.PHONY: ci
ci: install setup lint build tests

