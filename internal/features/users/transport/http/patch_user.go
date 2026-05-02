package http

import (
	"fmt"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/request"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/types"
)

// PatchUserRequest represents the incoming JSON body for a partial user update (DTO).
type PatchUserRequest struct {
	FullName    types.Nullable[string] `json:"full_name"`
	PhoneNumber types.Nullable[string] `json:"phone_number"`
}

var phoneRegex = regexp.MustCompile(`^\+[0-9]+$`)

// Validate performs early HTTP-level validation on the incoming payload.
// It ensures that string lengths and formats (like phone numbers) are correct
// before passing the data down to the domain layer.
func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := utf8.RuneCountInString(*r.FullName.Value)
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := utf8.RuneCountInString(*r.PhoneNumber.Value)
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !phoneRegex.MatchString(*r.PhoneNumber.Value) {
				return fmt.Errorf("invalid `PhoneNumber` format")
			}
		}
	}

	return nil
}

// PatchUserResponse represents the outgoing JSON body after a partial user update (JSON).
type PatchUserResponse UserDTOResponse

// PatchUser processes the HTTP request to partially update an existing user
// with the given id, validating the incoming JSON body.
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
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

	var req PatchUserRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	userPatch := userPatchFromRequest(req)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
