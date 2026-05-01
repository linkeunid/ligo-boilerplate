package auth

import (
	"github.com/linkeunid/ligo"
)

// Module returns the auth module with auth service and guards.
func Module() ligo.Module {
	return ligo.NewModule("auth",
		// Providers - exported for use in other modules
		ligo.Providers(
			ligo.Export(ligo.Factory[*AuthService](NewAuthService)),
		),
	)
}
