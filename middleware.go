package pluto

import (
	"context"
	"net/http"

	"bitbucket.org/aukbit/pluto/server/router"
)

// serviceContextMiddleware Middleware that adds service instance
// available in handlers context
func serviceContextMiddleware(s *service) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// get context
			ctx := r.Context()
			// Note: service instance is always available in handlers context
			// under the general name > pluto
			ctx = context.WithValue(ctx, "pluto", s)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}