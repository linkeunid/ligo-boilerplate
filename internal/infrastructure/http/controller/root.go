package controller

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/config"
)

// RootController handles root API info requests.
type RootController struct {
	cfg *config.Config
	log ligo.Logger
}

// NewRootController creates a new root controller.
func NewRootController(cfg *config.Config, log ligo.Logger) *RootController {
	return &RootController{cfg: cfg, log: log}
}

// Initialize is called when the root module initializes.
func (c *RootController) Initialize() error {
	c.log.Info("Root controller initializing")
	return nil
}

// Ready is called after all modules initialize, before serving.
func (c *RootController) Ready() error {
	c.log.Info("Root controller ready")
	return nil
}

// Draining is called before shutdown begins.
func (c *RootController) Draining() error {
	c.log.Info("Root controller draining")
	return nil
}

// Shutdown is called during shutdown.
func (c *RootController) Shutdown() error {
	c.log.Info("Root controller shutting down")
	return nil
}

// Register implements the Registerable interface for compile-time safe hook registration.
func (c *RootController) Register(registry *ligo.HookRegistry) {
	registry.OnInit(c.Initialize)
	registry.OnBootstrap(c.Ready)
	registry.BeforeShutdown(c.Draining)
	registry.OnShutdown(c.Shutdown)
}

// Routes registers all routes for the root controller.
func (c *RootController) Routes(r ligo.Router) {
	r.Handle("GET", "/", c.Info)
}

// Info handles GET /
func (c *RootController) Info(ctx ligo.Context) error {
	c.log.Info("API info requested")
	return ctx.JSON(200, map[string]any{
		"name":        "Ligo Boilerplate",
		"version":     c.cfg.Version,
		"description": "Clean Architecture example using Ligo framework",
		"architecture": map[string]string{
			"pattern": "Clean Architecture + Repository Pattern",
			"layers":  "domain → usecase → infrastructure → cmd",
		},
		"endpoints": map[string]string{
			"GET    /":             "API info",
			"GET    /health":       "Health check",
			"GET    /users":        "List all users (public)",
			"GET    /users/:id":    "Get user by ID (requires auth)",
			"POST   /users":        "Create user (requires auth)",
			"PUT    /users/:id":    "Update user (requires auth)",
			"DELETE /users/:id":    "Delete user (requires admin)",
			"POST   /files/upload": "Upload file",
			"GET    /files/:id":    "Download file",
			"GET    /files":        "List files",
			"DELETE /files/:id":    "Delete file",
		},
		"examples": map[string]string{
			"list_users":  `curl http://localhost:8080/users`,
			"get_user":    `curl -H "Authorization: Bearer user:secret" http://localhost:8080/users/<id>`,
			"create_user": `curl -X POST -H "Authorization: Bearer user:secret" -H "Content-Type: application/json" -d '{"name":"Alice","email":"alice@example.com"}' http://localhost:8080/users`,
			"update_user": `curl -X PUT -H "Authorization: Bearer user:secret" -H "Content-Type: application/json" -d '{"name":"Alice Updated"}' http://localhost:8080/users/<id>`,
			"delete_user": `curl -X DELETE -H "Authorization: Bearer admin:secret" http://localhost:8080/users/<id>`,
		},
	})
}
