# Story 1.4: Frontend Project Scaffolding

**Epic:** 1 - Project Foundation & Design System
**Status:** in-progress
**Assigned:** Sally (Frontend Developer)

## User Story

As a **developer**,
I want **a properly structured Angular 21 frontend project**,
So that **I can build features using standalone components and Signals**.

## Acceptance Criteria

- [ ] Angular 21 project created with standalone components
- [ ] `src/app/features/` for feature modules exists
- [ ] `src/app/shared/components/` for shared components exists
- [ ] `src/app/shared/services/` for shared services exists
- [ ] `src/app/core/` for singleton services exists
- [ ] `src/app/layouts/` for layout components exists
- [ ] Tailwind CSS configured with custom theme colors
- [ ] ESLint and Prettier configured for code quality
- [ ] Environment files for dev/staging/prod exist
- [ ] Proxy configuration for API calls is set up
- [ ] Path aliases (@app, @shared, @features, @core) configured

## Technical Requirements

### Project Structure
```
msls-frontend/
├── src/
│   ├── app/
│   │   ├── core/
│   │   │   ├── services/
│   │   │   │   └── .gitkeep
│   │   │   ├── guards/
│   │   │   │   └── .gitkeep
│   │   │   ├── interceptors/
│   │   │   │   └── .gitkeep
│   │   │   └── models/
│   │   │       └── .gitkeep
│   │   ├── shared/
│   │   │   ├── components/
│   │   │   │   └── .gitkeep
│   │   │   ├── directives/
│   │   │   │   └── .gitkeep
│   │   │   └── pipes/
│   │   │       └── .gitkeep
│   │   ├── features/
│   │   │   └── .gitkeep
│   │   ├── layouts/
│   │   │   └── .gitkeep
│   │   ├── app.component.ts
│   │   ├── app.config.ts
│   │   └── app.routes.ts
│   ├── assets/
│   ├── environments/
│   │   ├── environment.ts
│   │   ├── environment.development.ts
│   │   └── environment.production.ts
│   └── styles/
│       ├── _variables.scss
│       └── styles.scss
├── tailwind.config.js
├── proxy.conf.json
├── .eslintrc.json
├── .prettierrc
├── tsconfig.json
└── package.json
```

### Tailwind Theme Colors
```javascript
colors: {
  primary: {
    50: '#eff6ff',
    100: '#dbeafe',
    200: '#bfdbfe',
    300: '#93c5fd',
    400: '#60a5fa',
    500: '#3b82f6',
    600: '#2563eb',
    700: '#1d4ed8',
    800: '#1e40af',
    900: '#1e3a8a',
  },
  // Secondary, success, warning, danger colors
}
```

### Path Aliases (tsconfig.json)
```json
{
  "compilerOptions": {
    "paths": {
      "@app/*": ["src/app/*"],
      "@core/*": ["src/app/core/*"],
      "@shared/*": ["src/app/shared/*"],
      "@features/*": ["src/app/features/*"],
      "@env/*": ["src/environments/*"]
    }
  }
}
```

### Proxy Configuration (proxy.conf.json)
```json
{
  "/api": {
    "target": "http://localhost:8080",
    "secure": false,
    "changeOrigin": true
  }
}
```

## Tasks

- [ ] 1. Create Angular 21 project with `ng new msls-frontend --style=scss --routing=true --ssr=false --standalone=true --strict=true`
- [ ] 2. Install and configure Tailwind CSS
- [ ] 3. Create custom color theme in tailwind.config.js
- [ ] 4. Create directory structure (core, shared, features, layouts)
- [ ] 5. Configure path aliases in tsconfig.json
- [ ] 6. Create environment files (dev, staging, prod)
- [ ] 7. Create proxy.conf.json for API calls
- [ ] 8. Install and configure ESLint with Angular rules
- [ ] 9. Install and configure Prettier
- [ ] 10. Create base styles.scss with Tailwind imports
- [ ] 11. Install Angular CDK for accessibility primitives
- [ ] 12. Create .gitignore for Angular projects
- [ ] 13. Update package.json scripts
- [ ] 14. Verify `npm start` works and shows default page

## Definition of Done

- [ ] All acceptance criteria met
- [ ] `npm run build` succeeds without errors
- [ ] `npm run lint` passes with no errors
- [ ] `npm start` launches dev server successfully
- [ ] Tailwind classes work in templates
- [ ] Path aliases resolve correctly
- [ ] Proxy forwards /api requests to backend
- [ ] Code follows Angular Style Guide
- [ ] All components use standalone: true
