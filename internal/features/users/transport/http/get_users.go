package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
)

// GetUsersResponse represents the outgoing JSON body after getting a user slice (JSON).
type GetUsersResponse []UserDTOResponse

// GetUsers processes the HTTP request to get a list of users according to the limit and offset parameters.
// It writes the JSON response.
func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, w)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit/'offset'")
		return
	}

	page := domain.NewPagination(limit, offset)
	userDomains, err := h.usersService.GetUsers(ctx, page)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		queryLimit  = "limit"
		queryOffset = "offset"
	)

	limit, errLimit := request.GetIntQueryParam(r, queryLimit)
	if errLimit != nil {
		errLimit = fmt.Errorf(
			"get `%s` query params: %w",
			queryLimit,
			errLimit,
		)
	}

	offset, errOffset := request.GetIntQueryParam(r, queryOffset)
	if errOffset != nil {
		errOffset = fmt.Errorf(
			"get `%s` query params: %w",
			queryOffset,
			errOffset,
		)
	}

	if errs := errors.Join(
		errLimit,
		errOffset,
	); errs != nil {
		return nil, nil, errs
	}

	return limit, offset, nil
}
