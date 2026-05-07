.PHONY: run build test lint migrate up down docker-build docker-run start stop

IMAGE_NAME=sre-bootcamp-one2n
IMAGE_TAG=$(shell cat VERSION)

up:
	docker compose up -d postgres

down:
	docker compose down

stop:
	docker compose stop

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker compose up -d api

migrate:
	@echo "Running migrations..."
	@docker compose exec postgres pg_isready -U pguser -q || (echo "Postgres is not ready" && exit 1)
	go run cmd/migrate/main.go

start: docker-build up
	@echo "Waiting for postgres to be healthy..."
	@until docker compose exec postgres pg_isready -U pguser -q; do sleep 1; done
	@echo "Running migrations..."
	go run cmd/migrate/main.go
	@echo "Starting API..."
	docker compose up -d api
	@echo "All services up. API available at http://localhost:8080"
