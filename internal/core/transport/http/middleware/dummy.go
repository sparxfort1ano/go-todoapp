package middleware

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
)

// DebugLogger is a development middleware that traces the execution flow.
// It logs a specific message before passing the request to the next handler,
// and logs again when the response comes back up the chain.
func DebugLogger(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)

			log.Debug(fmt.Sprintf("--> before %s", s))

			next.ServeHTTP(w, r)

			log.Debug(fmt.Sprintf("<-- after %s", s))
		})
	}
}
