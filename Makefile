MAKEFLAGS += --silent

deps:
	@GOBIN=${PWD}/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

fmt:
	@gofmt -s -w .

lint:
	@./bin/golangci-lint run

test:
	@go test -race -cover -tags test ./...

cover:
	@go test -covermode=count -coverprofile cover.out ./...
	@go tool cover -html=cover.out -o cover.html
	@open cover.html

clean:
	@rm bin/*
	@go clean

gencerts:
	pushd examples/ssl && ./gencerts.sh && popd

run:
	go run ./examples/main.go

runssl: gencerts
	go run ./examples/main.go -cert examples/ssl/bundle.crt -key examples/ssl/server.key
