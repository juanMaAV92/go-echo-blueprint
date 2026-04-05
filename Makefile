.PHONY: run build test lint tidy

run:
	go run main.go

build:
	go build -o bin/app main.go

test:
	go test ./... -race -coverprofile=coverage.out

lint:
	golangci-lint run --timeout=5m

tidy:
	go mod tidy
