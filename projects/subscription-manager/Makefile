.PHONY: dev dev-backend dev-frontend build up down logs clean

# Development
dev:
	@echo "Starting development environment..."
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL..."
	@sleep 3
	@make dev-backend &
	@make dev-frontend

dev-backend:
	cd backend && go run ./cmd/api

dev-frontend:
	cd frontend && npm run dev

# Docker
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

# Database
db-shell:
	docker-compose exec postgres psql -U postgres -d subtrack

# Clean
clean:
	docker-compose down -v
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/node_modules

# Install dependencies
install:
	cd frontend && npm install
	cd backend && go mod download

# Testing
test:
	cd backend && go test -v ./...

# Lint
lint:
	cd frontend && npm run lint
	cd backend && golangci-lint run
