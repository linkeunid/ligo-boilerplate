# Modules Guide

## Creating a New Module

### 1. Create Package Directory

```bash
mkdir internal/product
```

### 2. Define Files

```
internal/product/
├── controller.go    # HTTP handlers
├── service.go       # Business logic
├── repository.go    # Data access
└── module.go        # Module definition
```

### 3. Implement Repository

```go
// repository.go
type ProductRepository struct {
    products map[string]*Product
}

func NewProductRepository() *ProductRepository {
    return &ProductRepository{
        products: make(map[string]*Product),
    }
}

func (r *ProductRepository) FindByID(id string) (*Product, bool) {
    p, found := r.products[id]
    return p, found
}

func (r *ProductRepository) Save(product *Product) {
    r.products[product.ID] = product
}
```

### 4. Implement Service

```go
// service.go
type ProductService struct {
    repo *ProductRepository
    log  ligo.Logger
}

func NewProductService(repo *ProductRepository, log ligo.Logger) *ProductService {
    return &ProductService{repo: repo, log: log}
}

func (s *ProductService) GetProduct(id string) (*Product, error) {
    product, found := s.repo.FindByID(id)
    if !found {
        return nil, common.ErrNotFound
    }
    return product, nil
}
```

### 5. Implement Controller

```go
// controller.go
type Controller struct {
    svc *ProductService
    log ligo.Logger
}

func NewController(svc *ProductService, log ligo.Logger) ligo.Controller {
    return &Controller{svc: svc, log: log}
}

func (c *Controller) Routes(r ligo.Router) {
    cr := ligo.NewChainRouter(r.Group("/products"))

    cr.GET("", c.ListProducts).
        Filter(common.GlobalExceptionFilter(c.log)).
        Intercept(common.LoggingInterceptor(c.log)).
        Handle()

    cr.GET("/:id", c.GetProduct).
        Filter(common.GlobalExceptionFilter(c.log)).
        Guard(auth.AuthGuard(authService)).
        Intercept(common.LoggingInterceptor(c.log)).
        Handle()
}

func (c *Controller) GetProduct(ctx ligo.Context) error {
    id := ctx.Param("id")
    product, err := c.svc.GetProduct(id)
    if err != nil {
        return err
    }
    return ctx.OK(product)
}
```

### 6. Define Module

```go
// module.go
func Module() ligo.Module {
    return ligo.NewModule("product",
        ligo.Imports(auth.Module()),
        ligo.Providers(
            ligo.Export(ligo.Factory[*ProductRepository](NewProductRepository)),
            ligo.Export(ligo.Factory[*ProductService](NewProductService)),
        ),
        ligo.Controllers(
            func(svc *ProductService, log ligo.Logger) ligo.Controller {
                return NewController(svc, log)
            },
        ),
    )
}
```

### 7. Register in Main

```go
app.Register(
    auth.Module(),
    product.Module(),
)
```

## Module Features

| Feature | Description | Example |
|---------|-------------|---------|
| **Guards** | Pre-handler authorization | `Guard(auth.AuthGuard())` |
| **Filters** | Exception handling | `Filter(common.GlobalExceptionFilter())` |
| **Interceptors** | Around-advice logic | `Intercept(common.LoggingInterceptor())` |
| **Pipes** | Request transformation | `Pipe(ValidationPipe())` |

## Existing Modules

| Module | Purpose | Endpoints |
|--------|---------|-----------|
| `auth` | Authentication/authorization | Guards, interceptors |
| `user` | User CRUD | `/users/*` |
| `file` | File upload/download | `/files/*` |
| `health` | Health check | `/health` |
| `root` | API info | `/` |
