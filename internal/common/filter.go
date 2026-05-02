package common

import (
	"errors"
	"fmt"

	"github.com/linkeunid/ligo"
)

// GlobalExceptionFilter handles all errors and converts them to HTTP responses.
func GlobalExceptionFilter(log ligo.Logger) ligo.ExceptionFilter {
	return func(err error, ctx ligo.Context) error {
		if err == nil {
			return nil
		}

		log.Error("Request error",
			ligo.LoggerField{Key: "method", Value: ctx.Request().Method},
			ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
			ligo.LoggerField{Key: "error", Value: err.Error()},
		)

		statusCode := 500
		message := "Internal Server Error"

		switch {
		case errors.Is(err, ErrUnauthorized):
			statusCode = 401
			message = "Unauthorized"
		case errors.Is(err, ErrForbidden):
			statusCode = 403
			message = "Forbidden"
		case errors.Is(err, ErrNotFound):
			statusCode = 404
			message = "Not Found"
		case errors.Is(err, ErrValidation):
			statusCode = 400
			message = "Validation Failed"
		}

		return ctx.JSON(statusCode, map[string]string{
			"error": message,
			"code":  fmt.Sprintf("%d", statusCode),
		})
	}
}
