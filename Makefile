# VERSION := $(shell git describe --tags)
# BUILD := $(shell git rev-parse --short HEAD)
# PROJECT := $(shell basename "$(PWD)")
# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
MAKEFLAGS += --silent

test:
	go test ./...

clean:
	rm bin/$(PROJECT)
	go clean

build:
	go build $(LDFLAGS) -o bin/$(PROJECT) ./examples/main.go

run:
	go run ./examples/main.go
