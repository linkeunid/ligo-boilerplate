package health

import (
	"github.com/linkeunid/ligo"
)

// Module returns the health module.
func Module() ligo.Module {
	return ligo.NewModule("health",
		ligo.Controllers(
			NewController,
		),
	)
}
