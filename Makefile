MAKEFLAGS += --silent

deps:
	@GOBIN=${PWD}/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

lint:
	@./bin/golangci-lint run

test:
	go test -race -cover -tags test ./...

cover:
	go test -covermode=count -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html
	open cover.html

clean:
	rm bin/*
	go clean

run:
	go run ./examples/main.go
