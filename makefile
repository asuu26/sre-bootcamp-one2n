.PHONY: run build test migrate up down docker-build docker-run

IMAGE_NAME=sre-bootcamp-one2n
IMAGE_TAG=0.1.0

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

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run:
	docker run --rm \
		--env-file .env \
		-e DB_HOST=postgres \
		--network sre-bootcamp-one2n_default \
		-p 8080:8080 \
		$(IMAGE_NAME):$(IMAGE_TAG)
