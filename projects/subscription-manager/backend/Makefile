.PHONY: run build test clean migrate-up migrate-down

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/

migrate-up:
	@echo "Migrations are run automatically on startup"

migrate-down:
	@echo "Run down migrations manually if needed"

tidy:
	go mod tidy

lint:
	golangci-lint run

dev:
	air
