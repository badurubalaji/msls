# Story 1.1: Backend Project Scaffolding

**Epic:** 1 - Project Foundation & Design System
**Status:** in-progress
**Assigned:** James (Backend Developer)

## User Story

As a **developer**,
I want **a properly structured Go backend project with all necessary dependencies and configurations**,
So that **I can start building features following established patterns**.

## Acceptance Criteria

- [ ] Go project structure matches architecture specifications
- [ ] `cmd/api/main.go` entry point exists
- [ ] `internal/modules/` directory structure is created
- [ ] `internal/pkg/` for shared utilities exists
- [ ] `configs/` with environment-based configuration
- [ ] `migrations/` directory for database migrations
- [ ] All core dependencies installed (Gin, GORM, sqlc, validator, uuid)
- [ ] Makefile includes targets for build, test, run, migrate
- [ ] Docker Compose file includes PostgreSQL 16, Redis, MinIO services
- [ ] `.env.example` documents all required environment variables

## Technical Requirements

### Project Structure
```
msls-backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── modules/
│   │   └── .gitkeep
│   ├── middleware/
│   │   └── .gitkeep
│   ├── pkg/
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   └── response/
│   └── config/
├── migrations/
├── api/
│   └── openapi.yaml
├── build/
│   └── docker/
│       ├── Dockerfile
│       └── docker-compose.yml
├── scripts/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Dependencies (go.mod)
```
github.com/gin-gonic/gin v1.9+
gorm.io/gorm v1.25+
gorm.io/driver/postgres v1.5+
github.com/redis/go-redis/v9 v9.0+
github.com/golang-jwt/jwt/v5 v5.0+
github.com/google/uuid v1.6+
github.com/go-playground/validator/v10 v10.0+
github.com/spf13/viper v1.18+
go.uber.org/zap v1.27+
golang.org/x/crypto v0.20+
```

### Docker Compose Services
- PostgreSQL 16 (port 5432)
- Redis 7 (port 6379)
- MinIO (ports 9000, 9001)

### Makefile Targets
- `make build` - Build the binary
- `make run` - Run the server
- `make test` - Run tests
- `make lint` - Run golangci-lint
- `make migrate-up` - Run migrations
- `make migrate-down` - Rollback migrations
- `make docker-up` - Start Docker services
- `make docker-down` - Stop Docker services

## Tasks

- [ ] 1. Initialize Go module with `go mod init github.com/anthropics/msls-backend`
- [ ] 2. Create directory structure as per architecture
- [ ] 3. Add all core dependencies to go.mod
- [ ] 4. Create `cmd/api/main.go` with basic Gin server
- [ ] 5. Create config package with Viper for environment loading
- [ ] 6. Create logger package with Zap structured logging
- [ ] 7. Create Dockerfile with multi-stage build
- [ ] 8. Create docker-compose.yml with PostgreSQL, Redis, MinIO
- [ ] 9. Create Makefile with all required targets
- [ ] 10. Create .env.example with all environment variables
- [ ] 11. Create .gitignore for Go projects
- [ ] 12. Create README.md with setup instructions
- [ ] 13. Verify `make docker-up && make run` works

## Definition of Done

- [ ] All acceptance criteria met
- [ ] `make build` succeeds without errors
- [ ] `make test` passes (even if no tests yet)
- [ ] `make lint` passes with no errors
- [ ] Docker services start successfully
- [ ] Server starts and responds to health check
- [ ] Code follows Google Go Style Guide
