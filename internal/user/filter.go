package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

// GlobalExceptionFilter handles all errors and converts them to HTTP responses.
func GlobalExceptionFilter(log ligo.Logger) ligo.ExceptionFilter {
	return func(err error, ctx ligo.Context) error {
		if err == nil {
			return nil
		}

		// Log the error
		log.Error("Request error",
			ligo.LoggerField{Key: "method", Value: ctx.Request().Method},
			ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
			ligo.LoggerField{Key: "error", Value: err.Error()},
		)

		// Convert to HTTP response based on error type
		statusCode := 500
		message := "Internal Server Error"

		switch {
		case errors.Is(err, common.ErrUnauthorized):
			statusCode = 401
			message = "Unauthorized"
		case errors.Is(err, common.ErrForbidden):
			statusCode = 403
			message = "Forbidden"
		case errors.Is(err, common.ErrNotFound):
			statusCode = 404
			message = "Not Found"
		case errors.Is(err, common.ErrValidation):
			statusCode = 400
			message = "Validation Failed"
		case strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "required"):
			statusCode = 400
			message = "Bad Request"
		}

		return ctx.JSON(statusCode, map[string]string{
			"error": message,
			"code":  fmt.Sprintf("%d", statusCode),
		})
	}
}

// AuthExceptionFilter handles authentication-related errors.
func AuthExceptionFilter() ligo.ExceptionFilter {
	return func(err error, ctx ligo.Context) error {
		if errors.Is(err, common.ErrUnauthorized) {
			return ctx.JSON(401, map[string]string{
				"error": "Unauthorized - Please provide a valid authentication token",
				"code":  "401",
			})
		}
		return err
	}
}

// ForbiddenExceptionFilter handles authorization-related errors.
func ForbiddenExceptionFilter() ligo.ExceptionFilter {
	return func(err error, ctx ligo.Context) error {
		if errors.Is(err, common.ErrForbidden) {
			return ctx.JSON(403, map[string]string{
				"error": "Forbidden - You don't have permission to access this resource",
				"code":  "403",
			})
		}
		return err
	}
}

// ValidationExceptionFilter handles validation-related errors.
func ValidationExceptionFilter() ligo.ExceptionFilter {
	return func(err error, ctx ligo.Context) error {
		if errors.Is(err, common.ErrValidation) {
			return ctx.JSON(400, map[string]any{
				"error":   "Validation failed",
				"code":    "400",
				"details": err.Error(),
			})
		}
		return err
	}
}
