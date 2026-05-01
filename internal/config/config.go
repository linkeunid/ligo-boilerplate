package config

import "github.com/linkeunid/ligo"

// Config holds application configuration
type Config struct {
	Env          string
	LogLevel     string
	MaxTodos     int
	DefaultPriority string
}

// NewConfig creates default configuration
func NewConfig() *Config {
	return &Config{
		Env:          "development",
		LogLevel:     "info",
		MaxTodos:     1000,
		DefaultPriority: "medium",
	}
}

// Module returns the config module with exported Config
func Module() ligo.Module {
	return ligo.NewModule("config",
		ligo.Providers(
			// Export makes Config available to other modules that import this one
			ligo.Export(ligo.Factory[*Config](NewConfig)),
		),
	)
}
