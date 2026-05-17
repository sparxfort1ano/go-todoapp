package http

import (
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// DeleteUser processes the HTTP request to delete all the info
// about the user with the given identificator.
//
// @Summary		Удаление пользователя
// @Description Удаление существующего в системе пользователя по его ID.
// @Tags 		users
// @Param 		id 			path int true 					"ID удаляемого пользователя"
// @Success 	204 										"Успешное удаление пользователя"
// @Failure 	400 		{object} response.ErrorResponse "Bad request"
// @Failure 	404 		{object} response.ErrorResponse "User not found"
// @Failure 	500 		{object} response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
