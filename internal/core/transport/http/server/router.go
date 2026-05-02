package server

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/middleware"
)

type APIVersion string

const (
	APIVersion1 = APIVersion("v1")
	APIVersion2 = APIVersion("v2")
	APIVersion3 = APIVersion("v3")
)

// APIVersionRouter represents a multiplexer for a specific API version (e.g., v1, v2).
// It encapsulates version-specific middlewares and a collection of routes.
//
// Example usage:
//
//	v2 := server.NewAPIVersionRouter(
//		server.APIVersion2,
//		middleware.DebugLogger("api v2 middleware"),
//	)
//	v2.RegisterRoutes(usersHandler.Routes()...)
type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	middleware []middleware.Middleware
}

// NewAPIVersionRouter creates a new instance of APIVersionRouter.
func NewAPIVersionRouter(
	apiVersion APIVersion,
	middleware ...middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

// RegisterRoutes binds individual API endpoints to this specific version router.
// It automatically wraps each handler with the route's specific middleware.
func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.withMiddleware())
	}
}

func (r *APIVersionRouter) withMiddleware() http.Handler {
	return middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
