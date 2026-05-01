package middleware

import (
	"fmt"
	"time"

	"github.com/linkeunid/ligo"
)

// RequestID adds a unique request ID to each request
func RequestID(next ligo.HandlerFunc) ligo.HandlerFunc {
	return func(ctx ligo.Context) error {
		// Generate simple request ID
		requestID := ctx.Request().Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		ctx.Set("request_id", requestID)
		ctx.Response().Header().Set("X-Request-ID", requestID)
		return next(ctx)
	}
}

// CORS adds basic CORS headers
func CORS(next ligo.HandlerFunc) ligo.HandlerFunc {
	return func(ctx ligo.Context) error {
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		// Handle preflight
		if ctx.Request().Method == "OPTIONS" {
			return ctx.String(200, "")
		}

		return next(ctx)
	}
}
