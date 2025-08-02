SHELL := /bin/bash
PATH := $(GOPATH)/bin:$(PATH)

.PHONY: all clean lint test cover bench

lint:
	golangci-lint run --fix

test:
	go test -race ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

bench: # includes tests
	go test -bench=. ./... -run=^$$

clean:
	go clean -testcache

all: clean lint test bench
