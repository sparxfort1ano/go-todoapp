// Package utils provides helper functions for parsing and
// validating HTTP request data.
package utils

import (
	"fmt"
	"net/http"
	"strconv"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// GetIntQueryParam extracts an integer query parameter
// from the HTTP request by its key.
// It returns nil if the parameter is missing
// or an error if the value is not a valid integer.
func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			errs.ErrInvalidArgument,
		)
	}

	return &val, nil
}
