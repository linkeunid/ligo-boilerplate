package module

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/config"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
)

// Health returns the health check module.
func Health() ligo.Module {
	return ligo.NewModule("health",
		ligo.Controllers(
			func(cfg *config.Config) ligo.Controller {
				return controller.NewHealthController(cfg)
			},
		),
	)
}
