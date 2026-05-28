.PHONY: run build tidy test

run:
	go run ./cmd/server

build:
	go build -o bin/tiktok-api ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...
