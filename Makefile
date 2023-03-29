PROJECT_NAME=jarvis
MODULE_NAME=github.com/kodep/jarvis

.DEFAULT_GOAL := build

.PHONY: build
build:
	@go build ./cmd/jarvis/

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: get
get:
	@go mod download

.PHONY: gen
gen:
	go generate ./...
