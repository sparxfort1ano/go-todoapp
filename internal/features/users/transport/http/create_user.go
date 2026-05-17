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
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164,min=10" example:"+79051707732"`
}

// CreateUserResponse represents the outgoing JSON body after a user is created (JSON).
type CreateUserResponse UserDTOResponse

// CreateUser processes the HTTP request to register a new user.
// It decodes the payload and writes the JSON response.
//
// @Summary		Создать пользователя
// @Description Создать нового пользователя в системе.
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		request 	body 		CreateUserRequest true 	"CreateUser тело запроса"
// @Success 	201 		{object} 	CreateUserResponse 		"Успешно созданный пользователь"
// @Failure 	400 		{object} 	response.ErrorResponse 	"Bad request"
// @Failure 	500 		{object} 	response.ErrorResponse 	"Internal server error"
// @Router 		/users 		[post]
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
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
