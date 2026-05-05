package module

import (
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/http/controller"
)

// Health returns the health check module.
func HealthModule() ligo.Module {
	return ligo.NewModule("health",
		// Use HookedController for compile-time safe hook registration.
		ligo.Controllers(ligo.HookedController(controller.NewHealthController)),
		// Module-level init hook: runs when health module initializes
		ligo.OnModuleInit(func() error {
			// Controller-level hooks handle logging
			return nil
		}),
		// Module-level destroy hook: runs when health module destroys
		ligo.OnModuleDestroy(func() error {
			// Controller-level hooks handle logging
			return nil
		}),
	)
}
