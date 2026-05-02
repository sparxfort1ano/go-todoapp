// Package http acts as the transport layer for the Users feature.
// It is responsible for parsing HTTP requests, formatting responses and routing.
package http

import (
	"context"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
)

// UsersHTTPHandler handles HTTP requests related to user management.
// All of its methods delegate the logic to the service layer.
type UsersHTTPHandler struct {
	usersService UsersService
}

// UsersService defines the contract that decouples the HTTP transport layer
// from the underlying domain logic.
type UsersService interface {
	// CreateUser enforces business rules (like length and symbol checks) on the user domain.
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	// GetUsers enforces business rules (like negative values in limit or offset parameter)
	// on the user domain.
	GetUsers(
		ctx context.Context,
		page domain.Pagination,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	// PatchUser at first requests to get the given user data by identificator,
	// then enforces business rules on both the user patch and user domain levels
	// (see ApplyPatch for details).
	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

// NewUsersHTTPHandler creates a new instance of UsersHTTPHandler.
func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

// Routes returns a list of HTTP routes to be registered in the server router.
//
// Example with route-specific middleware:
//
//		{
//		Method:     http.MethodGet,
//		Path:       "/users",
//		Handler:    h.GetUsers,
//		Middleware: []middleware.Middleware{
//			middleware.DebugLogger("get users middleware"),
//		},
//	}
func (h *UsersHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
