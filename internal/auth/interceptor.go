package auth

import (
	"github.com/linkeunid/ligo"
)

// AuditInterceptor logs admin actions for security auditing.
func AuditInterceptor(log ligo.Logger) ligo.Interceptor {
	return func(ctx ligo.Context, next ligo.HandlerFunc) error {
		user, _ := ctx.Get(ContextKeyUser).(*User)

		err := next(ctx)

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
