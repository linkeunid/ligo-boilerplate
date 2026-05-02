package middleware

import "github.com/linkeunid/ligo"

// BindPathParams populates _ligo_bound_data with the named path parameters.
// Must be added via .Use() before .Pipe() on a route so pipes can read the values.
//
// Example:
//
//	cr.GET("/:id", h).Use(BindPathParams("id")).Pipe(ligo.UUIDPipe("id")).Handle()
func BindPathParams(names ...string) ligo.Middleware {
	return func(next ligo.HandlerFunc) ligo.HandlerFunc {
		return func(ctx ligo.Context) error {
			params := make(map[string]string, len(names))
			for _, name := range names {
				params[name] = ctx.Param(name)
			}
			ctx.Set("_ligo_bound_data", params)
			return next(ctx)
		}
	}
}
