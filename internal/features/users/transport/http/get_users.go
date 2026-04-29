package http

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/utils"
)

// GetUsersResponse represents the outgoing JSON body after getting a user slice (JSON).
type GetUsersResponse []UserDTOResponse

// GetUsers processes the HTTP request to get a list of users according to the limit and offset parameters.
// It writes the JSON response.
func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetUsers handler")

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit/'offset'")
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

// GetIntQueryParam extracts an integer query parameter
// from the HTTP request by the keys offset and limit.
// It returns nil if the parameter is missing
// or an error if the value is not a valid integer.
func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get `limit` query params: %w", err)
	}

	offset, err := utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get `offset` query params: %w", err)
	}

	return limit, offset, nil
}
