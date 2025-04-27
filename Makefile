SHELL := bash
GO ?= go
TIMEOUT_UNIT ?= 5m

.PHONY: all
all: test

.PHONY: test
test:
	$(GO) test -race -timeout $(TIMEOUT_UNIT) ./...

.PHONY: test-unit-coverage
test-unit-coverage:
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic -timeout $(TIMEOUT_UNIT) ./...

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  test               - Run unit tests"
	@echo "  test-unit-coverage - Run unit tests with coverage reporting"
