# Modules Guide

## Creating a New Module

Follow Clean Architecture: define the domain first, then work outward.

### 1. Domain — entity and repository interface

```go
// internal/domain/entity/product.go
package entity

type Product struct {
    ID    string
    Name  string
    Price float64
}
```

```go
// internal/domain/repository/product.go
package repository

import "github.com/linkeunid/ligo-boilerplate/internal/domain/entity"

type ProductRepository interface {
    FindByID(id string) (*entity.Product, bool)
    FindAll() []*entity.Product
    Create(name string, price float64) *entity.Product
    Delete(id string) bool
}
```

### 2. UseCase — business logic

```go
// internal/usecase/product.go
package usecase

import (
    "github.com/linkeunid/ligo"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

type ProductUseCase struct {
    repo repository.ProductRepository
    log  ligo.Logger
}

func NewProductUseCase(repo repository.ProductRepository, log ligo.Logger) *ProductUseCase {
    return &ProductUseCase{repo: repo, log: log}
}

func (uc *ProductUseCase) GetProduct(id string) (*entity.Product, error) {
    p, found := uc.repo.FindByID(id)
    if !found {
        return nil, ErrNotFound
    }
    return p, nil
}
```

### 3. Persistence — repository implementation

```go
// internal/infrastructure/persistence/memory/product_repo.go
package memory

import (
    "sync"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

type ProductRepository struct {
    mu       sync.RWMutex
    products map[string]*entity.Product
}

func NewProductRepository() repository.ProductRepository {
    return &ProductRepository{products: make(map[string]*entity.Product)}
}

func (r *ProductRepository) FindByID(id string) (*entity.Product, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    p, found := r.products[id]
    return p, found
}

// ... implement remaining interface methods
```

### 4. Controller — HTTP handler

```go
// internal/infrastructure/http/controller/product.go
package controller

import (
    "github.com/linkeunid/ligo"
    "github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

type ProductController struct {
    uc  *usecase.ProductUseCase
    log ligo.Logger
}

func NewProductController(uc *usecase.ProductUseCase, log ligo.Logger) *ProductController {
    return &ProductController{uc: uc, log: log}
}

func (c *ProductController) Routes(r ligo.Router) {
    cr := ligo.NewChainRouter(r.Group("/products"))
    cr.GET("", c.ListProducts).Handle()
    cr.GET("/:id", c.GetProduct).Pipe(ligo.UUIDPipe("id")).Handle()
}

func (c *ProductController) GetProduct(ctx ligo.Context) error {
    id := ctx.Param("id") // validated UUID string, stored by UUIDPipe
    product, err := c.uc.GetProduct(id)
    if err != nil {
        return err
    }
    return ctx.OK(product)
}
```

### 5. Module — wire the layers

```go
// internal/module/product.go
package module

import (
    "github.com/linkeunid/ligo"
    "github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
    "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
    "github.com/linkeunid/ligo-boilerplate/internal/infrastructure/persistence/memory"
    "github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

func Product() ligo.Module {
    return ligo.NewModule("product",
        ligo.Providers(
            ligo.Factory[repository.ProductRepository](memory.NewProductRepository),
            ligo.Factory[*usecase.ProductUseCase](usecase.NewProductUseCase),
        ),
        ligo.Controllers(controller.NewProductController),
    )
}
```

`ligo.Controllers()` accepts the constructor directly. The DI binder resolves its parameters (`*usecase.ProductUseCase`, `ligo.Logger`) from the container at startup.

### 6. Register in main

```go
// cmd/api/main.go
app.Register(
    module.Product(),
    // ...
)
```

## Adding Auth Guards

Guards and middleware are typically wired directly inside the controller constructor and applied per-route using `ligo.NewChainRouter`:

```go
// internal/infrastructure/http/controller/product.go
func NewProductController(uc *usecase.ProductUseCase, jwt *infraauth.JWTAuth, log ligo.Logger) *ProductController {
    return &ProductController{
        uc:          uc,
        log:         log,
        authGuard:   infraauth.AuthGuard(jwt),
        exceptionMW: middleware.ExceptionMiddleware(log),
        loggingMW:   middleware.LoggingMiddleware(log),
    }
}

func (c *ProductController) Routes(r ligo.Router) {
    cr := ligo.NewChainRouter(r.Group("/products"))
    cr.Use(c.exceptionMW, c.loggingMW)

    cr.GET("", c.GetAll).Handle()
    cr.GET("/:id", c.GetByID).Guard(c.authGuard).Pipe(ligo.UUIDPipe("id")).Handle()
    cr.POST("", c.Create).Guard(c.authGuard).Pipe(ligo.ValidationPipe(&dto.CreateProductInput{})).Handle()
    cr.DELETE("/:id", c.Delete).Guard(c.authGuard).Pipe(ligo.UUIDPipe("id")).Handle()
}
```

The module registers `*infraauth.JWTAuth` as a provider so it is injected into the controller constructor:

```go
func Product() ligo.Module {
    return ligo.NewModule("product",
        ligo.Providers(
            ligo.Factory[repository.ProductRepository](memory.NewProductRepository),
            ligo.Factory[*usecase.ProductUseCase](usecase.NewProductUseCase),
            // JWTAuth is provided by the auth module or registered here
        ),
        ligo.Controllers(controller.NewProductController),
    )
}
```

## Existing Modules

| Module | Function | Endpoints |
|--------|---------|-----------|
| `auth` | `module.Auth()` | Provides `*infraauth.JWTAuth` |
| `user` | `module.User()` | `/users/*` |
| `file` | `module.File()` | `/files/*` |
| `health` | `module.Health()` | `/health` |
| `root` | `module.Root()` | `/` |

## Swapping Persistence

The repository pattern means replacing the storage backend requires only changing the concrete implementation in the module — the domain and usecase layers are untouched:

```go
// From in-memory:
repo := memory.NewProductRepository()

// To PostgreSQL (when implemented):
repo := postgres.NewProductRepository(db)
```
