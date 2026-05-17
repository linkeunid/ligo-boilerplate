package module

import (
	"github.com/linkeunid/ligo"

	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
)

// RootModule returns the root API info module.
func RootModule() ligo.Module {
	return ligo.NewModule(
		"root",
		// Use HookedController for compile-time safe hook registration.
		ligo.Controllers(ligo.HookedController(controller.NewRootController)),
		// Module-level init hook: runs when root module initializes
		ligo.OnModuleInit(func() error {
			// Controller-level hooks handle logging
			return nil
		}),
		// Module-level destroy hook: runs when root module destroys
		ligo.OnModuleDestroy(func() error {
			// Controller-level hooks handle logging
			return nil
		}),
	)
}
