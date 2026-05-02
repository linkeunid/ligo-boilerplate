package middleware

import (
	"time"

	"github.com/linkeunid/ligo"
)

// LoggingMiddleware logs request details.
func LoggingMiddleware(log ligo.Logger) ligo.Middleware {
	return func(next ligo.HandlerFunc) ligo.HandlerFunc {
		return func(ctx ligo.Context) error {
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
}
