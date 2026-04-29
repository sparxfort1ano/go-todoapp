package server

import (
	"fmt"
	"net/http"
)

type APIVersion string

const (
	APIVersion1 = APIVersion("v1")
	APIVersion2 = APIVersion("v2")
	APIVersion3 = APIVersion("v3")
)

// APIVersionRouter defines an API version to specific handler multiplexer.
type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
}

// NewAPIVersionRouter creates a new instance of APIVersionRouter.
func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

// RegisterRoutes binds individual API endpoints to the specific sub-router.
// It maps HTTP methods and endpoints to their handlers.
func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		// "METHOD /endpoint"
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
