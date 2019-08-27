all: test build

test:
	@go test -v ./...
.PHONY: test

deps:
	@go get -u && go mod tidy
.PHONY: deps

build:
	@go build -v ./...
.PHONY: build
