package user

import (
	"time"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/auth"
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

// AuditInterceptor logs admin actions for security auditing.
func AuditInterceptor(log ligo.Logger) ligo.Interceptor {
	return func(ctx ligo.Context, next ligo.HandlerFunc) error {
		// Get current user if available
		user, _ := ctx.Get(auth.ContextKeyUser).(*auth.User)

		err := next(ctx)

		// Log admin actions
		if user != nil && user.IsAdmin() {
			log.Info("Admin action performed",
				ligo.LoggerField{Key: "admin_id", Value: user.ID},
				ligo.LoggerField{Key: "action", Value: ctx.Request().Method},
				ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
				ligo.LoggerField{Key: "success", Value: err == nil},
			)
		}

		return err
	}
}
