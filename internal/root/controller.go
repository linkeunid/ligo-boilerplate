package root

import (
	"github.com/linkeunid/ligo"
)

// Controller handles root API info endpoint.
type Controller struct {
	version string
}

// NewController creates a new root controller.
func NewController() ligo.Controller {
	return &Controller{version: "0.7.0"}
}

// Routes registers all routes for the root module.
func (c *Controller) Routes(r ligo.Router) {
	r.Handle("GET", "/", c.Info)
}

// Info handles GET / - returns API information
func (c *Controller) Info(ctx ligo.Context) error {
	return ctx.JSON(200, map[string]any{
		"name":        "Ligo Boilerplate Example",
		"version":     c.version,
		"description": "Example application demonstrating all Ligo features",
		"features": []string{
			"Modules",
			"Dependency Injection",
			"Controllers",
			"Guards (Authorization)",
			"Pipes (Validation)",
			"Interceptors (Logging, Auditing)",
			"Exception Filters",
			"Middleware",
			"Lifecycle Hooks",
		},
		"endpoints": map[string]string{
			"GET    /":          "API info",
			"GET    /health":     "Health check",
			"GET    /users":      "List all users (public)",
			"GET    /users/:id":  "Get user by ID (requires auth)",
			"POST   /users":      "Create user (requires auth)",
			"PUT    /users/:id":  "Update user (requires auth)",
			"DELETE /users/:id":  "Delete user (requires admin)",
		},
		"authentication": map[string]any{
			"format":  "Bearer <token>",
			"example": "Authorization: Bearer user:secret",
			"tokens": map[string]string{
				"user":  "user:secret  (regular user)",
				"admin": "admin:secret (admin user)",
			},
		},
		"curl_examples": []string{
			"curl http://localhost:8080/",
			"curl http://localhost:8080/users",
			"curl -H 'Authorization: Bearer user:secret' http://localhost:8080/users/1",
			"curl -X POST -H 'Authorization: Bearer user:secret' -H 'Content-Type: application/json' -d '{\"name\":\"John\",\"email\":\"john@example.com\"}' http://localhost:8080/users",
			"curl -X DELETE -H 'Authorization: Bearer admin:secret' http://localhost:8080/users/1",
		},
	})
}
