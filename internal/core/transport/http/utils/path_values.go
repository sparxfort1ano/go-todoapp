package utils

import (
	"fmt"
	"net/http"
	"strconv"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// GetIntPathValue extracts an integer endpoint parameter
// from the HTTP request by its key.
// It returns an error if the value is not a valid integer.
func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			errs.ErrInvalidArgument,
		)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path value='%s' by key='%s' not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			errs.ErrInvalidArgument,
		)
	}

	return val, nil
}
