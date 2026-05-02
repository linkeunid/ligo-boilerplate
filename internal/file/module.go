package file

import (
	"github.com/linkeunid/ligo"
)

func Module() ligo.Module {
	return ligo.NewModule("file",
		ligo.Providers(
			ligo.Factory[*Controller](NewController),
		),
		ligo.Controllers(
			func(ctrl *Controller) ligo.Controller {
				return ctrl
			},
		),
	)
}
