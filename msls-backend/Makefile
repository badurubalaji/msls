# MSLS Backend Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=msls-api
BINARY_UNIX=$(BINARY_NAME)_unix

# Directories
CMD_DIR=./cmd/api
BUILD_DIR=./bin
DOCKER_DIR=./build/docker

# Environment
-include .env
export

.PHONY: all build clean run test lint fmt vet tidy deps docker-build docker-up docker-down \
	migrate-up migrate-down migrate-down-all migrate-status migrate-goto migrate-create migrate-force migrate-reset seed help \
	docs install-swagger

# Default target
all: lint test build

## Build targets

# Build the binary
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(CMD_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCMD) clean

## Run targets

# Run the application
run:
	@echo "Running..."
	$(GOCMD) run $(CMD_DIR)/main.go

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
run-dev:
	@echo "Running with hot reload..."
	air

## Test targets

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests in short mode
test-short:
	@echo "Running short tests..."
	$(GOTEST) -v -short ./...

## Code quality targets

# Run linter (requires golangci-lint)
lint:
	@echo "Linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 0; \
	fi

# Format code
fmt:
	@echo "Formatting..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Vetting..."
	$(GOVET) ./...

## Dependency targets

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify

## Docker targets

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t msls-backend:latest -f $(DOCKER_DIR)/Dockerfile .

# Start Docker services
docker-up:
	@echo "Starting Docker services..."
	docker compose -f $(DOCKER_DIR)/docker-compose.yml up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	docker compose -f $(DOCKER_DIR)/docker-compose.yml down

# View Docker logs
docker-logs:
	@echo "Viewing Docker logs..."
	docker compose -f $(DOCKER_DIR)/docker-compose.yml logs -f

# Remove Docker volumes
docker-clean:
	@echo "Removing Docker volumes..."
	docker compose -f $(DOCKER_DIR)/docker-compose.yml down -v

## Database migration targets

# Migration directories
MIGRATIONS_DIR=./migrations

# Run all migrations up
migrate-up:
	@echo "Running all migrations up..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up; \
	else \
		echo "golang-migrate not installed. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Run migrations down (1 step)
migrate-down:
	@echo "Running migrations down (1 step)..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1; \
	else \
		echo "golang-migrate not installed. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Run all migrations down (rollback everything)
migrate-down-all:
	@echo "Running all migrations down..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all; \
	else \
		echo "golang-migrate not installed. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Show migration version/status
migrate-status:
	@echo "Checking migration status..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version; \
	else \
		echo "golang-migrate not installed. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Migrate to a specific version
migrate-goto:
	@read -p "Enter version to migrate to: " version; \
	if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" goto $$version; \
	else \
		echo "golang-migrate not installed. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	fi

# Create new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# Force migration version (use with caution - fixes dirty state)
migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $$version

# Drop everything and re-run all migrations (development only)
migrate-reset: migrate-down-all migrate-up
	@echo "Migration reset complete."

# Seed the database with test data
seed:
	@echo "Seeding database..."
	$(GOCMD) run ./cmd/seed/main.go

## Utility targets

# Generate code (if using go generate)
generate:
	@echo "Generating code..."
	$(GOCMD) generate ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## Documentation targets

# Install swagger tool
install-swagger:
	@echo "Installing swag..."
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
docs:
	@echo "Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/api/main.go -o ./docs; \
	else \
		echo "swag not installed. Run: make install-swagger"; \
		exit 1; \
	fi

# Check environment
check-env:
	@echo "Checking environment..."
	@echo "Go version: $$(go version)"
	@echo "Docker version: $$(docker --version)"
	@echo "Docker Compose version: $$(docker compose version)"

# Help target
help:
	@echo "MSLS Backend Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Build targets:"
	@echo "  build         Build the binary"
	@echo "  build-linux   Build for Linux"
	@echo "  clean         Clean build artifacts"
	@echo ""
	@echo "Run targets:"
	@echo "  run           Run the application"
	@echo "  run-dev       Run with hot reload (requires air)"
	@echo ""
	@echo "Test targets:"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  test-short    Run short tests"
	@echo ""
	@echo "Code quality targets:"
	@echo "  lint          Run linter"
	@echo "  fmt           Format code"
	@echo "  vet           Run go vet"
	@echo ""
	@echo "Dependency targets:"
	@echo "  tidy          Tidy dependencies"
	@echo "  deps          Download dependencies"
	@echo "  verify        Verify dependencies"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-up     Start Docker services"
	@echo "  docker-down   Stop Docker services"
	@echo "  docker-logs   View Docker logs"
	@echo "  docker-clean  Remove Docker volumes"
	@echo ""
	@echo "Migration targets:"
	@echo "  migrate-up      Run all migrations up"
	@echo "  migrate-down    Run migrations down (1 step)"
	@echo "  migrate-down-all Roll back all migrations"
	@echo "  migrate-status  Show current migration version"
	@echo "  migrate-goto    Migrate to a specific version"
	@echo "  migrate-create  Create new migration files"
	@echo "  migrate-force   Force set migration version (fix dirty state)"
	@echo "  migrate-reset   Reset database (down-all + up)"
	@echo "  seed            Seed database with test data"
	@echo ""
	@echo "Documentation targets:"
	@echo "  docs          Generate Swagger documentation"
	@echo "  install-swagger Install swagger tool"
	@echo ""
	@echo "Utility targets:"
	@echo "  install-tools Install development tools"
	@echo "  check-env     Check environment"
	@echo "  help          Show this help message"
