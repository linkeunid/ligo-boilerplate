package root

import (
	"github.com/linkeunid/ligo"
)

// Module returns the root module with API info endpoint.
func Module() ligo.Module {
	return ligo.NewModule("root",
		ligo.Controllers(
			NewController,
		),
	)
}
