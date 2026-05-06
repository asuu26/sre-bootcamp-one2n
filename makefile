.PHONY: run build test migrate up down

up:
	docker compose up -d

down:
	docker compose down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test ./...

migrate:
	go run cmd/migrate/main.go
