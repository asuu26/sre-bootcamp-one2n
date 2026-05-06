# SRE Bootcamp - Student REST API

A CRUD REST API for managing student records, built with Go and Gin.

## Tech Stack

- **Language:** Go 1.26
- **Framework:** Gin
- **Database:** PostgreSQL 16
- **Logging:** Uber Zap (structured JSON)
- **Migrations:** golang-migrate
- **Container:** Docker + Docker Compose

## Project Structure

```
.
├── cmd/
│   ├── api/            # API server entrypoint
│   └── migrate/        # DB migration entrypoint
├── internal/
│   ├── api/handlers/   # HTTP handlers
│   ├── db/             # DB connection
│   ├── logger/         # Zap logger + Gin middleware
│   └── models/         # Request/response structs
├── migrations/         # SQL migration files
├── docker-compose.yml
├── Makefile
├── postman_collection.json
└── .env.example
```

## Prerequisites

- Docker + Docker Compose
- GNU Make

## One-Click Local Setup

**1. Clone the repo**
```bash
git clone https://github.com/75asu/sre-bootcamp-one2n.git
cd sre-bootcamp-one2n
```

**2. Copy env file**
```bash
cp .env.example .env
```

**3. Start everything**
```bash
make start
```

That's it. `make start` will:
1. Build the Docker image
2. Start Postgres and wait until healthy
3. Run DB migrations
4. Start the API container

Server runs at `http://localhost:8080`.

**Stop all services**
```bash
make down
```

## Local Development (without Docker)

Requires Go 1.26+ installed locally.

```bash
make up       # start Postgres
make migrate  # run migrations
make run      # run API locally
```

## Environment Variables

| Variable     | Description              | Default   |
|--------------|--------------------------|-----------|
| PORT         | Server port              | 8080      |
| GIN_MODE     | Gin mode (debug/release) | debug     |
| DB_HOST      | Postgres host            | localhost |
| DB_PORT      | Postgres port            | 5432      |
| DB_USER      | Postgres user            | pguser    |
| DB_PASSWORD  | Postgres password        | pgpass    |
| DB_NAME      | Postgres database        | sre_bootcamp |
| DB_SSLMODE   | SSL mode                 | disable   |

## API Endpoints

| Method | Endpoint                  | Description       |
|--------|---------------------------|-------------------|
| GET    | /healthcheck              | Health check      |
| POST   | /api/v1/students          | Create student    |
| GET    | /api/v1/students          | Get all students  |
| GET    | /api/v1/students/:id      | Get student by ID |
| PUT    | /api/v1/students/:id      | Update student    |
| DELETE | /api/v1/students/:id      | Delete student    |

### Example Request

```bash
curl -s -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","age":22}' | jq
```

### Example Response

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "age": 22,
  "created_at": "2026-05-06T11:34:00Z",
  "updated_at": "2026-05-06T11:34:00Z"
}
```

## Running with Docker

**Build the image**
```bash
make docker-build
```

**Run the container**
```bash
make up
make docker-run
```

The container connects to Postgres via the Docker Compose network. `DB_HOST` is overridden to `postgres` (the compose service name) at runtime.

## Available Make Commands

| Command           | Description                        |
|-------------------|------------------------------------|
| make up           | Start Docker containers            |
| make down         | Stop Docker containers             |
| make run          | Run the API server locally         |
| make build        | Build the binary                   |
| make test         | Run unit tests                     |
| make migrate      | Run DB migrations                  |
| make docker-build | Build the Docker image             |
| make docker-run   | Run the API in a Docker container  |
| make start        | One-click: build, migrate, run all |
| make stop         | Stop all containers                |

## Testing

```bash
make test
```

Import `postman_collection.json` into Postman to test all endpoints manually.
