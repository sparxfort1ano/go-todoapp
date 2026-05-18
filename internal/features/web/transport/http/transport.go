// Package http acts as the transport layer for the Web feature.
// It is responsible for transfering HTML files, formatting responses and routing.
package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
)

// WebHTTPHandler handles HTTP responses related to web management.
type WebHTTPHandler struct {
	webService WebService
}

// WebService defines the contract that decouples the HTTP transport layer
// from the underlying domain logic.
type WebService interface {
	// GetMainPage specifies the location of the index.html file and returns its data.
	GetMainPage() ([]byte, error)
}

// NewWebHTTPHandler creates a new instance of WebHTTPHandler.
func NewWebHTTPHandler(
	webService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

// Routes returns a list of HTTP routes to be registered in the server router.
func (h *WebHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{$}",
			Handler: h.GetMainPage,
		},
	}
}
