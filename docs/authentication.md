# Authentication

## Overview

Authentication is implemented as a domain service (`domain/service.AuthService`) with a JWT-style implementation in `infrastructure/auth/`. Guards enforce access control at the route level.

## How It Works

### 1. AuthService interface (domain layer)

```go
// internal/domain/service/auth.go
type AuthService interface {
    ValidateToken(token string) (User, error)
}
```

The domain defines the contract. The implementation lives in infrastructure.

### 2. JWTAuth implementation

```go
// internal/infrastructure/auth/jwt.go
type JWTAuth struct { log ligo.Logger }

func (j *JWTAuth) ValidateToken(token string) (service.User, error) {
    // For demo: accepts "role:secret" format tokens
    // Replace with real JWT parsing in production
}
```

### 3. Guards

```go
// internal/infrastructure/auth/guard.go

// AuthGuard validates the Bearer token and stores the user in context.
func AuthGuard(authService service.AuthService) ligo.Guard

// AdminGuard checks that the context user has admin role.
func AdminGuard() ligo.Guard
```

## Token Format (Demo)

```
Authorization: Bearer <role>:secret
```

| Token | Role |
|-------|------|
| `user:secret` | Regular user |
| `admin:secret` | Admin |

```bash
curl -H "Authorization: Bearer user:secret" http://localhost:8080/users/1
curl -H "Authorization: Bearer admin:secret" -X DELETE http://localhost:8080/users/1
```

## Using Guards in Routes

```go
import (
    infraauth "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/auth"
)

authSvc := infraauth.NewJWTAuth(log)
authGuard := infraauth.AuthGuard(authSvc)
adminGuard := infraauth.AdminGuard()

cr.GET("/:id", c.GetUser).Guard(authGuard).Handle()
cr.DELETE("/:id", c.DeleteUser).Guard(authGuard, adminGuard).Handle()
```

## Reading the Authenticated User

```go
import "github.com/linkeunid/ligo-boilerplate/internal/domain/service"

func (c *Controller) Handler(ctx ligo.Context) error {
    user, ok := ctx.Get(service.ContextKeyUser).(service.User)
    if !ok {
        return usecase.ErrUnauthorized
    }

    if user.IsAdmin() {
        // admin-only logic
    }

    return ctx.OK(map[string]string{"id": user.GetID()})
}
```

## Error Types

Errors are defined in `internal/usecase/errors.go`:

```go
var (
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrNotFound     = errors.New("resource not found")
    ErrValidation   = errors.New("validation failed")
)
```

`ExceptionMiddleware` maps these and additional framework errors to HTTP status codes:

| Error | Status |
|-------|--------|
| `usecase.ErrUnauthorized` | 401 |
| `usecase.ErrForbidden` | 403 |
| `usecase.ErrNotFound` | 404 |
| `ligo.ErrBadRequest` | 400 |
| `validator.ValidationErrors` | 422 |
| `usecase.ErrValidation` | 400 |
| everything else | 500 |

`ligo.ErrBadRequest` is wrapped by `UUIDPipe`, `ParseIntPipe`, and `ParseBoolPipe` when a path parameter is invalid. `validator.ValidationErrors` is returned by `ValidationPipe` when struct tag validation fails.

## Audit Logging

Admin actions are automatically logged by `middleware.AuditMiddleware`:

```go
middleware.AuditMiddleware(log)

// Logs: Admin action performed | admin_id=1 | action=DELETE | path=/users/2 | success=true
```

Apply it per-route:

```go
cr.DELETE("/:id", c.Delete).
    Guard(authGuard, adminGuard).
    Use(middleware.AuditMiddleware(log)).
    Handle()
```

## Production Recommendations

### Replace with real JWT

```go
import "github.com/golang-jwt/jwt/v5"

func (j *JWTAuth) ValidateToken(token string) (service.User, error) {
    t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return []byte(j.secret), nil
    })
    if err != nil || !t.Valid {
        return nil, usecase.ErrUnauthorized
    }

    claims := t.Claims.(jwt.MapClaims)
    return &entity.User{
        ID:   claims["sub"].(string),
        Role: claims["role"].(string),
    }, nil
}
```

### Use environment variables for secrets

```go
cfg.JWTSecret = os.Getenv("JWT_SECRET")
if cfg.JWTSecret == "" {
    log.Fatal("JWT_SECRET is required")
}
```
