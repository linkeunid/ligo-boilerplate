package module

import (
	"github.com/linkeunid/ligo"

	"github.com/linkeunid/ligo-boilerplate/internal/infrastructure/worker"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

func WorkerModule() ligo.Module {
	return ligo.NewModule(
		"worker",
		ligo.Providers(
			ligo.Factory[*usecase.WorkerUseCase](usecase.NewWorkerUseCase),
		),
		ligo.Controllers(ligo.HookedController(worker.NewController)),
	)
}
