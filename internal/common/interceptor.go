package common

import (
	"time"

	"github.com/linkeunid/ligo"
)

// LoggingInterceptor logs request details.
func LoggingInterceptor(log ligo.Logger) ligo.Interceptor {
	return func(ctx ligo.Context, next ligo.HandlerFunc) error {
		start := time.Now()

		log.Debug("Request started",
			ligo.LoggerField{Key: "method", Value: ctx.Request().Method},
			ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
		)

		err := next(ctx)

		duration := time.Since(start)

		if err != nil {
			log.Error("Request completed with error",
				ligo.LoggerField{Key: "method", Value: ctx.Request().Method},
				ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
				ligo.LoggerField{Key: "duration_ms", Value: duration.Milliseconds()},
				ligo.LoggerField{Key: "error", Value: err.Error()},
			)
		} else {
			log.Debug("Request completed",
				ligo.LoggerField{Key: "method", Value: ctx.Request().Method},
				ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
				ligo.LoggerField{Key: "duration_ms", Value: duration.Milliseconds()},
			)
		}

		return err
	}
}
