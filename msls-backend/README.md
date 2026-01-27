# MSLS Backend

Multi-School Learning System (MSLS) Backend API built with Go.

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- **ORM**: [GORM](https://gorm.io/) - Go ORM library
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Object Storage**: MinIO
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Logging**: [Zap](https://github.com/uber-go/zap)
- **Authentication**: JWT with [golang-jwt](https://github.com/golang-jwt/jwt)

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose
- Make

## Quick Start

### 1. Clone and Setup

```bash
# Navigate to the backend directory
cd msls-backend

# Copy environment file
cp .env.example .env

# Download dependencies
make deps
```

### 2. Start Infrastructure

Start PostgreSQL, Redis, and MinIO:

```bash
make docker-up
```

### 3. Run the Server

```bash
make run
```

The API server will start at `http://localhost:8080`.

### 4. Verify Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T00:00:00Z"
  }
}
```

## Project Structure

```
msls-backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── modules/             # Feature modules (auth, users, schools, etc.)
│   ├── middleware/          # HTTP middleware
│   └── pkg/
│       ├── config/          # Configuration management
│       ├── database/        # Database connection
│       ├── logger/          # Structured logging
│       └── response/        # Standard API responses
├── migrations/              # Database migrations
├── api/                     # API specifications (OpenAPI)
├── build/
│   └── docker/
│       ├── Dockerfile
│       └── docker-compose.yml
├── scripts/                 # Utility scripts
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Available Make Commands

### Build and Run

| Command | Description |
|---------|-------------|
| `make build` | Build the binary |
| `make run` | Run the application |
| `make run-dev` | Run with hot reload (requires air) |
| `make clean` | Clean build artifacts |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make test-short` | Run short tests only |

### Code Quality

| Command | Description |
|---------|-------------|
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make vet` | Run go vet |

### Docker

| Command | Description |
|---------|-------------|
| `make docker-up` | Start Docker services |
| `make docker-down` | Stop Docker services |
| `make docker-logs` | View Docker logs |
| `make docker-clean` | Remove Docker volumes |

### Database Migrations

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback last migration |
| `make migrate-create` | Create a new migration |

### Utilities

| Command | Description |
|---------|-------------|
| `make deps` | Download dependencies |
| `make tidy` | Tidy go.mod |
| `make install-tools` | Install development tools |
| `make help` | Show all available commands |

### Documentation

| Command | Description |
|---------|-------------|
| `make docs` | Generate Swagger documentation |
| `make install-swagger` | Install swagger tool |

## Configuration

Configuration is loaded from environment variables. See `.env.example` for all available options.

### Key Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/staging/production) | development |
| `SERVER_PORT` | Server port | 8080 |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `REDIS_HOST` | Redis host | localhost |
| `JWT_SECRET` | JWT signing secret | (must be set) |

## API Documentation

This project uses Swagger/OpenAPI for API documentation.

### Generating Documentation

```bash
make docs
```

### Viewing Documentation

Once the server is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## API Endpoints

### Health Check

- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check endpoint

### API v1

- `GET /api/v1/ping` - Ping endpoint

### Documentation

- `GET /swagger/*` - Swagger UI and API documentation

## Development

### Install Development Tools

```bash
make install-tools
```

This installs:
- `golangci-lint` - Linter
- `air` - Hot reload
- `migrate` - Database migrations

### Code Style

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go/best-practices).

### Running Linter

```bash
make lint
```

### Hot Reload

For development with hot reload:

```bash
make run-dev
```

## Docker Services

| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5432 | Primary database |
| Redis | 6379 | Cache and sessions |
| MinIO | 9000 | Object storage API |
| MinIO Console | 9001 | MinIO web console |

### MinIO Console

Access the MinIO console at `http://localhost:9001`:
- Username: `minioadmin`
- Password: `minioadmin`

## License

Proprietary - All rights reserved.
