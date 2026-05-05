package controller

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/config"
)

// HealthController handles health check requests.
type HealthController struct {
	cfg *config.Config
	log ligo.Logger
}

// NewHealthController creates a new health controller.
func NewHealthController(cfg *config.Config, log ligo.Logger) *HealthController {
	return &HealthController{cfg: cfg, log: log}
}

// Initialize is called when the health module initializes.
func (c *HealthController) Initialize() error {
	c.log.Info("Health controller initializing")
	return nil
}

// Ready is called after all modules initialize, before serving.
func (c *HealthController) Ready() error {
	c.log.Info("Health controller ready")
	return nil
}

// Register implements the Registerable interface for compile-time safe hook registration.
func (c *HealthController) Register(registry *ligo.HookRegistry) {
	registry.OnInit(c.Initialize)
	registry.OnBootstrap(c.Ready)
}

// Routes registers all routes for the health controller.
func (c *HealthController) Routes(r ligo.Router) {
	r.Handle("GET", "/health", c.Check)
}

// Check handles GET /health
func (c *HealthController) Check(ctx ligo.Context) error {
	return ctx.JSON(200, map[string]string{
		"status":  "ok",
		"version": c.cfg.Version,
	})
}
