package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linkeunid/ligo"

	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

type Controller struct {
	workerUseCase *usecase.WorkerUseCase
	log           ligo.Logger
	cancel        context.CancelFunc
	mu            sync.Mutex
	running       atomic.Bool
	wg            sync.WaitGroup
}

func NewController(wuc *usecase.WorkerUseCase, log ligo.Logger) *Controller {
	return &Controller{
		workerUseCase: wuc,
		log:           log,
	}
}

func (c *Controller) Initialize() error {
	c.log.Info("Worker controller initializing")
	return nil
}

func (c *Controller) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running.Load() {
		return nil
	}

	c.log.Info("Worker controller starting worker")

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.running.Store(true)

	c.wg.Add(1)
	go c.run(ctx)
	return nil
}

func (c *Controller) Drain() error {
	c.log.Info("Worker controller draining - waiting for current work to complete")
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	c.log.Info("Worker controller drained")
	return nil
}

func (c *Controller) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running.Load() {
		return nil
	}

	c.log.Info("Worker controller stopping worker")
	c.running.Store(false)
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	c.log.Info("Worker controller stopped")
	return nil
}

func (c *Controller) Register(registry *ligo.HookRegistry) {
	registry.OnInit(c.Initialize)
	registry.OnBootstrap(c.Start)
	registry.BeforeShutdown(c.Drain)
	registry.OnShutdown(c.Stop)
}

func (c *Controller) run(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	c.workerUseCase.Execute()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("Worker controller: worker stopped")
			return
		case <-ticker.C:
			c.workerUseCase.Execute()
		}
	}
}
