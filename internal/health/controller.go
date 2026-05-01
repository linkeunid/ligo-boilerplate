package health

import (
	"github.com/linkeunid/ligo"
)

// Controller handles health check endpoints.
type Controller struct {
	version string
}

// NewController creates a new health controller.
func NewController() ligo.Controller {
	return &Controller{version: "0.7.0"}
}

// Routes registers all routes for the health module.
func (c *Controller) Routes(r ligo.Router) {
	r.Handle("GET", "/health", c.Check)
}

// Check handles GET /health
func (c *Controller) Check(ctx ligo.Context) error {
	return ctx.JSON(200, map[string]string{
		"status":  "ok",
		"version": c.version,
	})
}
