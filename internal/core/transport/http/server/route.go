package server

import "net/http"

// Route binds an HTTP method and URI pattern to specific handler.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// NewRoute creates a new instance of Route.
func NewRoute(method string, path string, handler http.HandlerFunc) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
