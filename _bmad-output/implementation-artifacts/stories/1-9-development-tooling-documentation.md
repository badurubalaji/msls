# Story 1.9: Development Tooling & Documentation

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** High

## User Story

As a **developer**,
I want **comprehensive development tooling and documentation**,
So that **new team members can onboard quickly and maintain code quality**.

## Acceptance Criteria

### Documentation
**Given** a new developer joins the project
**When** they read the documentation
**Then** README includes setup instructions for both backend and frontend
**And** Architecture decision records (ADRs) document key decisions
**And** API documentation is auto-generated (Swagger/OpenAPI)

### Code Quality Tooling
**Given** code is being developed
**When** pre-commit hooks run
**Then** Go code is formatted with gofmt and linted with golangci-lint
**Then** TypeScript code is formatted with Prettier and linted with ESLint
**And** Commit messages follow conventional commits format

### CI Pipeline
**Given** the CI pipeline runs
**When** a PR is submitted
**Then** All tests must pass
**And** Linting must pass with no errors
**And** Build must succeed for both backend and frontend

## Technical Requirements

### Backend Documentation

#### Swagger/OpenAPI Setup
```go
// Install swag CLI
// go install github.com/swaggo/swag/cmd/swag@latest

// Add annotations to handlers
// @Summary Get user by ID
// @Description Get user details by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Success{data=User}
// @Failure 404 {object} response.Error
// @Router /users/{id} [get]
```

#### Makefile Targets
```makefile
docs:
	swag init -g cmd/api/main.go -o ./docs

lint:
	golangci-lint run ./...

test:
	go test -v -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
```

### Frontend Documentation

#### ESLint Configuration
Already configured by Angular CLI, enhance with:
- @typescript-eslint rules
- Angular-specific rules
- Prettier integration

#### Prettier Configuration
Update `.prettierrc` if needed for team preferences.

### Pre-commit Hooks (Husky)

#### Backend (.pre-commit-config.yaml or Makefile)
```bash
# Before commit
make fmt
make lint
make test
```

#### Frontend (husky + lint-staged)
```json
// package.json
{
  "husky": {
    "hooks": {
      "pre-commit": "lint-staged",
      "commit-msg": "commitlint -E HUSKY_GIT_PARAMS"
    }
  },
  "lint-staged": {
    "*.{ts,html}": ["eslint --fix", "prettier --write"],
    "*.scss": ["prettier --write"]
  }
}
```

### Conventional Commits
```
feat: add user registration endpoint
fix: resolve null pointer in auth middleware
docs: update API documentation
style: format code with prettier
refactor: extract validation logic
test: add unit tests for user service
chore: update dependencies
```

### README Structure
```markdown
# MSLS - Multi-School Learning System

## Quick Start
### Prerequisites
### Backend Setup
### Frontend Setup

## Architecture
### Backend Structure
### Frontend Structure

## Development
### Running Tests
### Code Style
### Making Commits

## API Documentation
### Accessing Swagger UI
### Authentication

## Deployment
### Docker
### Environment Variables
```

## Tasks

### Backend
1. [ ] Install and configure swaggo for OpenAPI
2. [ ] Add Swagger annotations to existing endpoints
3. [ ] Set up Swagger UI endpoint
4. [ ] Configure golangci-lint with .golangci.yml
5. [ ] Update Makefile with lint, test, docs targets
6. [ ] Create README.md for backend

### Frontend
7. [ ] Configure ESLint with Angular rules
8. [ ] Ensure Prettier is configured
9. [ ] Install husky for git hooks
10. [ ] Configure lint-staged
11. [ ] Install commitlint for conventional commits
12. [ ] Update README.md for frontend

### Project-wide
13. [ ] Create main project README.md
14. [ ] Document environment variables
15. [ ] Create CONTRIBUTING.md guidelines
16. [ ] Create architecture decision record template
17. [ ] Document first ADR (tech stack selection)

## Definition of Done

- [ ] Swagger UI accessible at /swagger/
- [ ] golangci-lint passes with no errors
- [ ] ESLint and Prettier configured in frontend
- [ ] Pre-commit hooks work for both projects
- [ ] Conventional commits enforced
- [ ] README files complete and accurate
- [ ] New developer can set up project following README
