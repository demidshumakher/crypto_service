# Crypto Service

## What Is This?

Crypto Service is a Go HTTP REST API for registering users, authenticating with JWT tokens, and tracking cryptocurrency prices. It stores users, tracked cryptocurrencies, and price history in PostgreSQL, while current market data is fetched from the CoinGecko API.

The service listens on port `8080` and is designed to run either through Docker Compose or directly with a locally reachable PostgreSQL database.

## Key Functionality

- Register and log in users with bcrypt-hashed passwords.
- Issue JWT tokens for authenticated API access.
- Add cryptocurrencies for tracking by symbol, such as `BTC` or `ETH`.
- Fetch current tracked cryptocurrency data.
- Refresh prices from CoinGecko.
- Read price history and calculated price statistics.
- Configure and trigger scheduled price refreshes.

Primary routes:

| Method | Path | Authentication | Purpose |
|--------|------|----------------|---------|
| `POST` | `/auth/register` | No | Register a user and return a JWT token. |
| `POST` | `/auth/login` | No | Authenticate a user and return a JWT token. |
| `GET` | `/crypto` | Yes | List tracked cryptocurrencies. |
| `POST` | `/crypto` | Yes | Add a cryptocurrency by symbol. |
| `GET` | `/crypto/{symbol}` | Yes | Get one tracked cryptocurrency. |
| `PUT` | `/crypto/{symbol}/refresh` | Yes | Refresh one cryptocurrency from CoinGecko. |
| `GET` | `/crypto/{symbol}/history` | Yes | Get stored price history. |
| `GET` | `/crypto/{symbol}/stats` | Yes | Get price statistics from history. |
| `DELETE` | `/crypto/{symbol}` | Yes | Route is present; see Additional Notes. |
| `GET` | `/schedule` | Yes | Read refresh schedule settings. |
| `PUT` | `/schedule` | Yes | Update refresh schedule settings. |
| `POST` | `/schedule/trigger` | Yes | Trigger a refresh of all tracked cryptocurrencies. |

Authenticated requests must include:

```http
Authorization: Bearer <token>
```

## Architecture & Technology

The project follows a layered Go structure: HTTP handlers under `internal/rest`, business logic under `internal/service`, PostgreSQL persistence under `internal/repository/postgresql`, domain models under `domain`, and reusable packages under `pkg`. The executable entry point is `app/cryptoserver.go`.

Core technologies are Go, the standard `net/http` server and `ServeMux`, PostgreSQL, `pgx` through `database/sql`, JWT authentication, bcrypt password hashing, Docker, Docker Compose, and the external CoinGecko API.

## Prerequisites

- Go `1.24.2` or newer, based on `go.mod`.
- Docker and Docker Compose, for the recommended local setup.
- PostgreSQL, if running the service without Docker Compose.
- Python 3 and `pip3`, if running the provided Python integration tests.
- Network access to `https://api.coingecko.com/api/v3`.

## Installation

Install Go dependencies:

```bash
go mod download
```

Install Python test dependencies:

```bash
make install
```

Or install them directly:

```bash
pip3 install -r requirements.txt
```

## Environment Variables

Service runtime variables:

| Name | Required | Description | Example |
|------|-----------|-------------|---------|
| `DB_HOST` | No | PostgreSQL host used by the Go service. Defaults to `localhost`. | `postgres` |
| `DB_PORT` | No | PostgreSQL port used by the Go service. Defaults to `5432`. When running the app on the host against the Compose database, use `5433`. | `5432` |
| `DB_USER` | No | PostgreSQL username used by the Go service. Defaults to `postgres`. | `postgres` |
| `DB_PASSWORD` | No | PostgreSQL password used by the Go service. Defaults to `postgres`. | `postgres` |
| `DB_NAME` | No | PostgreSQL database name used by the Go service. Defaults to `cryptodb`. | `cryptodb` |
| `DB_SSL_MODE` | No | PostgreSQL SSL mode used in the connection string. Defaults to `disable`. | `disable` |

Docker Compose PostgreSQL variables:

| Name | Required | Description | Example |
|------|-----------|-------------|---------|
| `POSTGRES_USER` | Yes in Compose | Username created by the PostgreSQL container. | `postgres` |
| `POSTGRES_PASSWORD` | Yes in Compose | Password created by the PostgreSQL container. | `postgres` |
| `POSTGRES_DB` | Yes in Compose | Database created by the PostgreSQL container. | `cryptodb` |

Test variables:

| Name | Required | Description | Example |
|------|-----------|-------------|---------|
| `SCHEDULE` | No | When set to `1`, `make test` runs the additional schedule endpoint tests. | `1` |

The CoinGecko API key, CoinGecko base URL, JWT signing secret, JWT expiration, and HTTP listen port are currently hard-coded in `app/cryptoserver.go`.

## Running Locally

Recommended Docker Compose startup:

```bash
docker compose up --build
```

This starts:

- the Go API on `http://localhost:8080`
- PostgreSQL 15, exposed on host port `5433`
- database initialization from `db/schemes.sql` and `db/trigger.sql`

Register a user:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"demo-password"}'
```

Log in:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"demo-password"}'
```

Add a cryptocurrency:

```bash
curl -X POST http://localhost:8080/crypto \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <token>' \
  -d '{"symbol":"BTC"}'
```

List tracked cryptocurrencies:

```bash
curl http://localhost:8080/crypto \
  -H 'Authorization: Bearer <token>'
```

Run the service directly against the Compose PostgreSQL container:

```bash
DB_HOST=localhost \
DB_PORT=5433 \
DB_USER=postgres \
DB_PASSWORD=postgres \
DB_NAME=cryptodb \
DB_SSL_MODE=disable \
go run ./app/cryptoserver.go
```

## Development

Run package compilation checks:

```bash
go test ./...
```

Run the provided integration test script after the service is running:

```bash
make test
```

Run the additional schedule tests:

```bash
make test SCHEDULE=1
```

Useful Make targets:

| Command | Purpose |
|---------|---------|
| `make install` | Install Python test dependencies from `requirements.txt`. |
| `make test` | Run the main integration tests. |
| `make test SCHEDULE=1` | Run main and schedule integration tests. |
| `make clean` | Remove local temporary test/build artifacts listed in the Makefile. |
| `make help` | Print available Make targets. |

## Build & Deployment

Build a local binary:

```bash
go build -o cryptoserver ./app/cryptoserver.go
```

Build the Docker image:

```bash
docker build -t crypto-service .
```

Run through Docker Compose:

```bash
docker compose up --build -d
```

Stop the Compose stack:

```bash
docker compose down
```

The Dockerfile uses a multi-stage build. It compiles the Go binary from `app/cryptoserver.go`, copies the binary and `db/schemes.sql` into a small Alpine runtime image, exposes port `8080`, and runs `./cryptoserver`.

## Project Structure

```text
.
|-- app/
|   `-- cryptoserver.go          # Application entry point and dependency wiring
|-- db/
|   |-- schemes.sql              # PostgreSQL tables
|   `-- trigger.sql              # Price history trimming trigger
|-- domain/                      # Domain models and errors
|-- internal/
|   |-- repository/              # Persistence helpers and PostgreSQL repositories
|   |-- rest/                    # HTTP routing, handlers, middleware, responses
|   `-- service/                 # Business logic
|-- pkg/
|   |-- coingecko/               # CoinGecko API client
|   |-- jwt/                     # JWT generation and validation
|   `-- trigger/                 # Scheduled background work
|-- tests/
|   `-- tests.py                 # Python integration tests
|-- Dockerfile
|-- docker-compose.yml
|-- Makefile
|-- go.mod
`-- requirements.txt
```

## Troubleshooting

If the service exits with a database connection error, confirm PostgreSQL is running and the `DB_*` variables match the database location. With Docker Compose, the app container uses `DB_HOST=postgres` and `DB_PORT=5432`; a host-run Go process should use `DB_HOST=localhost` and `DB_PORT=5433`.

If adding or refreshing a cryptocurrency fails, confirm the service can reach CoinGecko and that the symbol exists in CoinGecko's `/coins/list` response.

If integration tests cannot reach the server, start the service first:

```bash
docker compose up -d
```

If schedule tests are skipped, run:

```bash
make test SCHEDULE=1
```

## Additional Notes

- The database stores users, cryptocurrencies, and prices. `db/trigger.sql` keeps at most 100 price rows per cryptocurrency.
- Price lookups use CoinGecko market data with `vs_currency=rub`.
- JWT tokens are signed with a hard-coded secret and use an expiration timestamp initialized at service startup.
- The HTTP server port is hard-coded to `8080`.
- Error responses are encoded with a `message` field by the current REST utility.
- The `DELETE /crypto/{symbol}` route is registered, but the current handler reads and returns history instead of deleting the record.
- The current README intentionally documents repository behavior as found in the code and configuration.
