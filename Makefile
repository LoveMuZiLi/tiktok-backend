.PHONY: run build tidy test

run:
	go run ./cmd/api

build:
	go build -o bin/tiktok-api ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...
