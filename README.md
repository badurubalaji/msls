# MSLS - Multi-School Learning System

A comprehensive multi-tenant learning management system designed to support multiple schools with isolated data and customizable features.

## Project Overview

MSLS (Multi-School Learning System) is a modern educational platform that enables:

- **Multi-tenancy**: Complete data isolation between schools using PostgreSQL Row-Level Security (RLS)
- **Role-based access**: Flexible permission system for administrators, teachers, students, and parents
- **Course management**: Create and manage courses, lessons, and educational content
- **Student tracking**: Monitor student progress, grades, and attendance
- **Communication**: Built-in messaging between teachers, students, and parents

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       Client Layer                          │
│                 Angular 21 SPA (Tailwind CSS)               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       API Gateway                           │
│                   Go + Gin Framework                        │
│           (JWT Auth, Rate Limiting, CORS)                   │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  PostgreSQL   │    │     Redis     │    │     MinIO     │
│    (RLS)      │    │    (Cache)    │    │   (Storage)   │
└───────────────┘    └───────────────┘    └───────────────┘
```

## Tech Stack

### Backend
- **Language**: Go 1.23+
- **Framework**: Gin HTTP framework
- **ORM**: GORM
- **Database**: PostgreSQL 16 with Row-Level Security
- **Cache**: Redis 7
- **Object Storage**: MinIO
- **Authentication**: JWT

### Frontend
- **Framework**: Angular 21
- **UI Styling**: Tailwind CSS 4
- **Testing**: Vitest
- **State Management**: Angular Signals

## Prerequisites

- Go 1.23 or later
- Node.js 20 or later
- npm 10 or later
- Docker and Docker Compose
- Make

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd msls
```

### 2. Backend Setup

```bash
cd msls-backend

# Copy environment configuration
cp .env.example .env

# Download dependencies
make deps

# Start infrastructure (PostgreSQL, Redis, MinIO)
make docker-up

# Run database migrations
make migrate-up

# Start the backend server
make run
```

The backend API will be available at `http://localhost:8080`.

### 3. Frontend Setup

```bash
cd msls-frontend

# Install dependencies
npm install

# Start development server
npm start
```

The frontend will be available at `http://localhost:4200`.

### 4. Verify Installation

```bash
# Check backend health
curl http://localhost:8080/health

# Open frontend in browser
open http://localhost:4200
```

## Project Structure

```
msls/
├── msls-backend/           # Go backend API
│   ├── cmd/api/            # Application entry point
│   ├── internal/           # Internal packages
│   │   ├── modules/        # Feature modules
│   │   ├── middleware/     # HTTP middleware
│   │   └── pkg/            # Shared packages
│   ├── migrations/         # Database migrations
│   ├── docs/               # Swagger documentation
│   └── build/docker/       # Docker configuration
│
├── msls-frontend/          # Angular frontend
│   ├── src/app/            # Application source
│   │   ├── core/           # Core services
│   │   ├── shared/         # Shared components
│   │   └── features/       # Feature modules
│   └── public/             # Static assets
│
├── docs/                   # Project documentation
│   └── adr/                # Architecture Decision Records
│
├── CONTRIBUTING.md         # Contribution guidelines
└── README.md               # This file
```

## Running Tests

### Backend Tests

```bash
cd msls-backend

# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Frontend Tests

```bash
cd msls-frontend

# Run tests in watch mode
npm test

# Run tests once
npm test -- --run

# Run tests with coverage
npm test -- --coverage
```

## API Documentation

Once the backend is running, access the Swagger UI documentation at:

```
http://localhost:8080/swagger/index.html
```

Generate or update API documentation:

```bash
cd msls-backend
make docs
```

## Development Tools

### Backend

```bash
cd msls-backend

# Install all development tools
make install-tools

# Run linter
make lint

# Format code
make fmt

# Run with hot reload
make run-dev
```

### Frontend

```bash
cd msls-frontend

# Run linter
npm run lint

# Fix linting issues
npm run lint:fix

# Format code
npm run format

# Check formatting
npm run format:check
```

## Environment Variables

### Backend (.env)

See `msls-backend/.env.example` for all available configuration options.

Key variables:
- `APP_ENV`: Environment (development/staging/production)
- `SERVER_PORT`: API server port (default: 8080)
- `DB_*`: Database connection settings
- `REDIS_*`: Redis connection settings
- `JWT_SECRET`: JWT signing secret (must be changed in production)

### Frontend

The frontend uses Angular environment files and a proxy configuration for API requests during development.

## Docker Services

| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5432 | Primary database |
| Redis | 6379 | Cache and sessions |
| MinIO API | 9000 | Object storage |
| MinIO Console | 9001 | MinIO web UI |

Start all services:

```bash
cd msls-backend
make docker-up
```

Stop all services:

```bash
make docker-down
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Architecture Decisions

Architecture Decision Records (ADRs) are stored in `docs/adr/`. See [ADR-001](docs/adr/001-tech-stack.md) for the technology stack decisions.

## License

Proprietary - All rights reserved.
