# MSLS Frontend

Multi-School Learning System (MSLS) Frontend built with Angular 21.

## Tech Stack

- **Framework**: [Angular 21](https://angular.dev/) - Modern web application framework
- **UI Styling**: [Tailwind CSS 4](https://tailwindcss.com/) - Utility-first CSS framework
- **Testing**: [Vitest](https://vitest.dev/) - Fast unit testing framework
- **Linting**: [ESLint](https://eslint.org/) with Angular and TypeScript rules
- **Formatting**: [Prettier](https://prettier.io/) - Code formatter
- **Git Hooks**: [Husky](https://typicode.github.io/husky/) with lint-staged and commitlint

## Prerequisites

- Node.js 20 or later
- npm 10 or later
- Angular CLI 21

## Quick Start

### 1. Install Dependencies

```bash
npm install
```

### 2. Start Development Server

```bash
npm start
```

The application will start at `http://localhost:4200` with API proxy configured to the backend.

### 3. Build for Production

```bash
npm run build:prod
```

## Available Scripts

| Command | Description |
|---------|-------------|
| `npm start` | Start development server with proxy |
| `npm run build` | Build the application |
| `npm run build:prod` | Production build with optimizations |
| `npm test` | Run unit tests with Vitest |
| `npm run lint` | Run ESLint |
| `npm run lint:fix` | Run ESLint with auto-fix |
| `npm run format` | Format code with Prettier |
| `npm run format:check` | Check code formatting |

## Project Structure

```
msls-frontend/
├── src/
│   ├── app/
│   │   ├── core/               # Core module (guards, interceptors, services)
│   │   ├── shared/             # Shared module (components, directives, pipes)
│   │   ├── features/           # Feature modules
│   │   └── app.component.ts    # Root component
│   ├── assets/                 # Static assets
│   ├── styles/                 # Global styles
│   └── main.ts                 # Application entry point
├── public/                     # Public static files
├── .husky/                     # Git hooks
├── angular.json                # Angular CLI configuration
├── eslint.config.js            # ESLint configuration
├── .prettierrc                 # Prettier configuration
├── .lintstagedrc               # lint-staged configuration
├── commitlint.config.js        # Commitlint configuration
├── tsconfig.json               # TypeScript configuration
└── package.json
```

## Development

### Code Style

This project enforces consistent code style using ESLint and Prettier:

- **ESLint**: Catches code quality issues and enforces Angular best practices
- **Prettier**: Ensures consistent code formatting

### Git Hooks

Pre-commit hooks run automatically to:

1. **Lint staged files**: ESLint fixes issues on staged TypeScript and HTML files
2. **Format code**: Prettier formats staged files
3. **Validate commits**: Commitlint enforces conventional commit messages

### Commit Message Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(auth): add login form validation
fix(dashboard): resolve chart rendering issue
docs(readme): update installation instructions
```

## API Proxy

Development server proxies API requests to the backend. Configuration is in `proxy.conf.json`:

```json
{
  "/api": {
    "target": "http://localhost:8080",
    "secure": false
  }
}
```

## Testing

Run tests with Vitest:

```bash
# Run tests in watch mode
npm test

# Run tests once
npm test -- --run

# Run tests with coverage
npm test -- --coverage
```

## Building

### Development Build

```bash
npm run build
```

### Production Build

```bash
npm run build:prod
```

Production builds include:
- AOT compilation
- Tree shaking
- Code splitting
- Minification
- Source maps (optional)

Output is placed in the `dist/` directory.

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

Proprietary - All rights reserved.
