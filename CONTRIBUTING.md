# Contributing to MSLS

Thank you for your interest in contributing to the Multi-School Learning System (MSLS). This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)

## Code of Conduct

- Be respectful and inclusive in all interactions
- Focus on constructive feedback
- Help others learn and grow
- Report unacceptable behavior to project maintainers

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Set up the development environment following the README.md instructions
4. Create a new branch for your feature or fix

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-fix-description
```

## Development Workflow

### Branch Naming

Use descriptive branch names with prefixes:

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications
- `chore/` - Maintenance tasks

Examples:
```
feature/user-authentication
fix/login-validation-error
docs/api-documentation
refactor/simplify-auth-middleware
```

### Development Process

1. Ensure your local main branch is up to date
2. Create a feature branch from main
3. Make your changes with appropriate tests
4. Run linters and tests locally
5. Commit using conventional commit messages
6. Push to your fork
7. Open a pull request

## Code Style

### Backend (Go)

This project follows the [Google Go Style Guide](https://google.github.io/styleguide/go/best-practices).

Key points:
- Use `gofmt` and `goimports` for formatting
- Follow naming conventions (exported names start with uppercase)
- Handle errors explicitly; avoid blank identifier for errors
- Write meaningful comments for exported types and functions
- Keep functions focused and reasonably sized

Run linting before committing:
```bash
cd msls-backend
make lint
make fmt
```

### Frontend (Angular/TypeScript)

Follow the [Angular Style Guide](https://angular.dev/style-guide).

Key points:
- Use single quotes for strings
- End statements with semicolons
- Use 2-space indentation
- Prefer `const` over `let` when possible
- Use meaningful component and service names

Run linting and formatting:
```bash
cd msls-frontend
npm run lint:fix
npm run format
```

## Commit Messages

This project uses [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation changes |
| `style` | Code style changes (formatting, semicolons, etc.) |
| `refactor` | Code refactoring without behavior change |
| `perf` | Performance improvements |
| `test` | Adding or updating tests |
| `build` | Build system or dependencies |
| `ci` | CI/CD changes |
| `chore` | Other maintenance tasks |
| `revert` | Revert a previous commit |

### Scope

Use the module or feature name as scope:
- `auth` - Authentication
- `users` - User management
- `schools` - School/tenant management
- `courses` - Course management
- `api` - API changes
- `ui` - UI components
- `config` - Configuration

### Examples

```
feat(auth): add JWT token refresh endpoint

fix(users): resolve duplicate email validation

docs(api): update swagger annotations for auth endpoints

refactor(middleware): simplify error handling logic

test(auth): add unit tests for password validation
```

### Breaking Changes

For breaking changes, add `BREAKING CHANGE:` in the footer:

```
feat(api): change response format for pagination

BREAKING CHANGE: Pagination now uses cursor-based approach instead of offset.
The `offset` parameter is replaced with `cursor`.
```

## Pull Request Process

### Before Opening a PR

1. Ensure all tests pass locally
2. Run linters and fix any issues
3. Update documentation if needed
4. Rebase on latest main if needed

### PR Requirements

- Clear, descriptive title following commit message format
- Description of changes and motivation
- Reference to related issues (if any)
- Screenshots for UI changes (if applicable)
- Tests for new functionality
- Documentation updates (if needed)

### PR Template

```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Related Issues
Fixes #123

## Testing
Describe how changes were tested.

## Checklist
- [ ] Tests pass locally
- [ ] Linting passes
- [ ] Documentation updated
- [ ] No new warnings
```

### Review Process

1. At least one approval required for merge
2. All CI checks must pass
3. Resolve all review comments
4. Squash commits if requested
5. Maintainers will merge approved PRs

## Testing Requirements

### Backend

- Write unit tests for all new functionality
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Mock external dependencies

Run tests:
```bash
cd msls-backend
make test
make test-coverage
```

### Frontend

- Write unit tests for components and services
- Test user interactions and edge cases
- Use Angular testing utilities

Run tests:
```bash
cd msls-frontend
npm test -- --run
npm test -- --coverage
```

### Coverage Expectations

- New features: Minimum 80% coverage
- Bug fixes: Tests that reproduce and verify the fix
- Refactoring: Maintain existing coverage

## Documentation

### Code Documentation

- Backend: Use Go doc comments for exported types and functions
- Frontend: Use JSDoc comments for public APIs

### API Documentation

- Update Swagger annotations for API changes
- Include request/response examples
- Document error responses

### README Updates

Update README files when:
- Adding new features
- Changing setup instructions
- Modifying available commands
- Adding new dependencies

## Questions?

If you have questions about contributing, please:
1. Check existing documentation
2. Search closed issues and PRs
3. Open a new issue for discussion

Thank you for contributing to MSLS!
