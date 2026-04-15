.PHONY: run test test-integration test-cover tidy lint build swagger

-include .env
export

run:
	go run cmd/main.go

build:
	go build -o main cmd/main.go

test:
	go test ./...

test-integration:
	go test -tags=integration ./tests/...

test-cover:
	go test ./... -coverpkg=./internal/service/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

tidy:
	go mod tidy

swagger:
	swag init -g cmd/main.go -o docs

lint:
	golangci-lint run --timeout=5m