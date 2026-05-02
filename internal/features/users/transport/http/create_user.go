package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// CreateUserRequest represents the incoming JSON body for creating a user (DTO).
type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number"  validate:"omitempty,e164,min=10"`
}

// CreateUserResponse represents the outgoing JSON body after a user is created (JSON).
type CreateUserResponse UserDTOResponse

// CreateUser processes the HTTP request to register a new user.
// It decodes the payload and writes the JSON response.
func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	var req CreateUserRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(req)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninialized(dto.FullName, dto.PhoneNumber)
}
