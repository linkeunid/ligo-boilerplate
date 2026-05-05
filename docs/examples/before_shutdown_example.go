package main

import (
	"time"

	"github.com/linkeunid/ligo"
)

// HTTPServer demonstrates the BeforeApplicationShutdown hook
// for graceful drain-stop scenarios.
type HTTPServer struct {
	addr string
	log  ligo.Logger
	// In production, you'd have an actual server here
}

func NewHTTPServer(addr string, log ligo.Logger) *HTTPServer {
	return &HTTPServer{addr: addr, log: log}
}

// OnModuleInit starts the HTTP server.
func (s *HTTPServer) OnModuleInit() error {
	s.log.Info("HTTP server starting", ligo.LoggerField{Key: "addr", Value: s.addr})
	// In production: start your HTTP server here
	return nil
}

// BeforeApplicationShutdown is called before shutdown begins.
// This is the ideal place to:
// 1. Stop accepting new connections
// 2. Wait for in-flight requests to complete (drain)
// 3. Prepare for shutdown
func (s *HTTPServer) BeforeApplicationShutdown() error {
	s.log.Info("HTTP server draining", ligo.LoggerField{Key: "addr", Value: s.addr})
	// In production:
	// 1. Set server to drain mode (stop accepting new connections)
	// 2. Wait for active requests to complete
	// 3. Return when ready to shutdown
	time.Sleep(100 * time.Millisecond) // Simulate drain time
	s.log.Info("HTTP server drained", ligo.LoggerField{Key: "addr", Value: s.addr})
	return nil
}

// OnApplicationShutdown is called after drain-stop.
// This is where you actually close the server.
func (s *HTTPServer) OnApplicationShutdown() error {
	s.log.Info("HTTP server shutting down", ligo.LoggerField{Key: "addr", Value: s.addr})
	// In production: close the server here
	return nil
}

// Example usage in a module:
//
//	func ServerModule() ligo.Module {
//	    return ligo.NewModule("server",
//	        ligo.Providers(
//	            ligo.Factory[*HTTPServer](NewHTTPServer),
//	        ),
//	    )
//	}
//
// The hooks will be called in order:
// 1. OnModuleInit - server starts
// 2. BeforeApplicationShutdown - server drains (no new connections, wait for in-flight)
// 3. OnApplicationShutdown - server closes
// 4. OnModuleDestroy - final cleanup

// Key differences between BeforeApplicationShutdown and OnApplicationShutdown:
//
// BeforeApplicationShutdown:
//   - Called FIRST during shutdown
//   - Ideal for drain-stop scenarios
//   - Stop accepting new work
//   - Wait for in-flight operations to complete
//   - Prepare for shutdown
//
// OnApplicationShutdown:
//   - Called AFTER drain-stop
//   - Actually close resources
//   - Final cleanup
//   - Release system resources
