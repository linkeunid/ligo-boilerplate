# Authentication

## Overview

This boilerplate uses a simple Bearer token authentication system. In production, replace this with JWT or OAuth2.

## How It Works

### 1. AuthGuard

Validates the `Authorization` header:

```go
Authorization: Bearer <token>
```

Format: `<username>:<password>`

Example tokens (for development):
- `user:secret` - Regular user
- `admin:secret` - Admin user

### 2. User Context

On successful auth, the user is stored in request context:

```go
user, _ := ctx.Get(auth.ContextKeyUser).(*auth.User)
```

### 3. Role-Based Access

Check user roles:

```go
if user.IsAdmin() {
    // Admin-only logic
}
```

Or use the built-in `RolesGuard`:

```go
cr.DELETE("/:id", c.Delete).
    Guard(auth.AuthGuard(authService), ligo.RolesGuard(auth.ContextKeyUser, "admin"))
```

## Protected Routes

```go
cr.GET("/users/:id", c.GetUser).
    Guard(auth.AuthGuard(authService))  // Requires auth
    Handle()

cr.DELETE("/users/:id", c.DeleteUser).
    Guard(auth.AuthGuard(authService), ligo.RolesGuard(auth.ContextKeyUser, "admin"))  // Requires admin
    Handle()
```

## User Model

```go
type User struct {
    ID    string
    Name  string
    Email string
    Role  string  // "user" or "admin"
}
```

## Production Recommendations

### 1. Replace with JWT

```go
import "github.com/golang-jwt/jwt/v5"

func JWTGuard(secret string) ligo.Guard {
    return func(ctx ligo.Context) (bool, error) {
        tokenString := strings.TrimPrefix(ctx.Request().Header.Get("Authorization"), "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            return false, common.ErrUnauthorized
        }

        claims := token.Claims.(jwt.MapClaims)
        user := &User{
            ID:   claims["sub"].(string),
            Role: claims["role"].(string),
        }
        ctx.Set(auth.ContextKeyUser, user)

        return true, nil
    }
}
```

### 2. Add Password Hashing

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 3. Use Environment Variables

```go
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
```

### 4. Add Token Refresh

Implement refresh token rotation for better security.

## Audit Trail

All admin actions are logged via `AuditInterceptor`:

```
Admin action performed | admin_id=admin123 | action=DELETE | path=/users/456 | success=true
```

Enable in your routes:

```go
Intercept(common.LoggingInterceptor(log), auth.AuditInterceptor(log))
```
