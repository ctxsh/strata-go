# VERSION := $(shell git describe --tags)
# BUILD := $(shell git rev-parse --short HEAD)
# PROJECT := $(shell basename "$(PWD)")
# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
MAKEFLAGS += --silent

test:
	go test -race -cover -tags test ./...

cover:
	go test -covermode=count -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html

clean:
	rm bin/$(PROJECT)
	go clean

examples:
	go build $(LDFLAGS) -o bin/$(PROJECT) ./examples/main.go

run:
	go run ./examples/main.go
