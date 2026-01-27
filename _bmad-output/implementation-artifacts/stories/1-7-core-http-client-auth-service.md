# Story 1.7: Core HTTP Client & Authentication Service

**Epic:** 1 - Project Foundation & Design System
**Status:** ready-for-dev
**Priority:** Critical

## User Story

As a **frontend developer**,
I want **a configured HTTP client with authentication interceptors**,
So that **all API calls automatically include tokens and handle auth errors**.

## Acceptance Criteria

**Given** a user is logged in
**When** an API request is made
**Then** the Authorization header includes the JWT access token
**And** the X-Tenant-ID header includes the current tenant ID

**Given** an API returns 401 Unauthorized
**When** a refresh token exists and is valid
**Then** the system automatically refreshes the access token
**And** the original request is retried with new token
**And** if refresh fails, user is redirected to login

**Given** the application needs to track requests
**When** any API call is made
**Then** a unique X-Request-ID header is included
**And** loading state is trackable via a global loading service

## Technical Requirements

### Service Location
`msls-frontend/src/app/core/`

### File Structure
```
core/
├── services/
│   ├── auth.service.ts        # Authentication state & methods
│   ├── api.service.ts         # HTTP client wrapper
│   ├── storage.service.ts     # Token storage (localStorage)
│   └── loading.service.ts     # Global loading state
├── interceptors/
│   ├── auth.interceptor.ts    # Adds Authorization header
│   ├── tenant.interceptor.ts  # Adds X-Tenant-ID header
│   ├── request-id.interceptor.ts # Adds X-Request-ID header
│   ├── error.interceptor.ts   # Handles HTTP errors
│   └── loading.interceptor.ts # Tracks loading state
├── guards/
│   ├── auth.guard.ts          # Route protection
│   └── tenant.guard.ts        # Tenant validation
└── models/
    ├── user.model.ts          # User interface
    ├── auth.model.ts          # Auth tokens interface
    └── api-response.model.ts  # API response wrapper
```

### Auth Service API
```typescript
@Injectable({ providedIn: 'root' })
export class AuthService {
  private currentUser = signal<User | null>(null);
  private accessToken = signal<string | null>(null);
  private refreshToken = signal<string | null>(null);

  readonly isAuthenticated = computed(() => !!this.accessToken());
  readonly user = this.currentUser.asReadonly();

  login(credentials: LoginRequest): Observable<AuthResponse>;
  logout(): void;
  refreshAccessToken(): Observable<string>;
  setTokens(tokens: AuthTokens): void;
  clearTokens(): void;
}
```

### Auth Interceptor
```typescript
export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const token = authService.accessToken();

  if (token) {
    req = req.clone({
      setHeaders: { Authorization: `Bearer ${token}` }
    });
  }

  return next(req).pipe(
    catchError(error => {
      if (error.status === 401) {
        return authService.refreshAccessToken().pipe(
          switchMap(newToken => {
            const retryReq = req.clone({
              setHeaders: { Authorization: `Bearer ${newToken}` }
            });
            return next(retryReq);
          }),
          catchError(() => {
            authService.logout();
            return throwError(() => error);
          })
        );
      }
      return throwError(() => error);
    })
  );
};
```

### Loading Service
```typescript
@Injectable({ providedIn: 'root' })
export class LoadingService {
  private activeRequests = signal(0);
  readonly isLoading = computed(() => this.activeRequests() > 0);

  show(): void { this.activeRequests.update(n => n + 1); }
  hide(): void { this.activeRequests.update(n => Math.max(0, n - 1)); }
}
```

## Tasks

1. [ ] Create user.model.ts with User interface
2. [ ] Create auth.model.ts with token interfaces
3. [ ] Create api-response.model.ts
4. [ ] Create storage.service.ts for token persistence
5. [ ] Create auth.service.ts with signal-based state
6. [ ] Create loading.service.ts with request counter
7. [ ] Create auth.interceptor.ts
8. [ ] Create tenant.interceptor.ts
9. [ ] Create request-id.interceptor.ts
10. [ ] Create error.interceptor.ts
11. [ ] Create loading.interceptor.ts
12. [ ] Create auth.guard.ts with canActivate
13. [ ] Create tenant.guard.ts
14. [ ] Register interceptors in app.config.ts
15. [ ] Write unit tests for services
16. [ ] Write unit tests for interceptors

## Definition of Done

- [ ] All services created using ng generate
- [ ] All interceptors implemented as functional interceptors
- [ ] Token refresh logic works correctly
- [ ] Loading state tracks active requests
- [ ] Guards protect routes properly
- [ ] Unit tests pass with >80% coverage
