package module

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/config"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
)

// Root returns the root API info module.
func Root() ligo.Module {
	return ligo.NewModule("root",
		ligo.Controllers(
			func(cfg *config.Config, log ligo.Logger) ligo.Controller {
				return controller.NewRootController(cfg, log)
			},
		),
	)
}
