// Package middleware provides HTTP interceptors that wrap standard handlers
// with common cross-cutting logic like logging, tracing and panic recovery.
package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// ChainMiddleware builds a single http.Handler from a chain of middleware functions.
// It applies the middleware in reverse order so they execute in the exact order provided.
func ChainMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
