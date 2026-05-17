package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// GetUserResponse represents the outgoing JSON body after a user is gotten (JSON).
type GetUserResponse UserDTOResponse

// GetUser processes the HTTP request to get a user with the given id.
// It writes the JSON response.
//
// @Summary		Получение пользователя
// @Description Получение конкретного пользователя по его ID.
// @Tags 		users
// @Produce 	json
// @Param 		id 	path int true 					"ID удаляемого пользователя"
// @Success 	200 {object} GetUserResponse 		"Успешное получение пользователя"
// @Failure 	400 {object} response.ErrorResponse "Bad request"
// @Failure 	404 {object} response.ErrorResponse "User not found"
// @Failure 	500 {object} response.ErrorResponse "Internal server error"
// @Router 		/users/{id}  [get]
func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	userID, err := request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JSONResponse(response, http.StatusOK)
}
