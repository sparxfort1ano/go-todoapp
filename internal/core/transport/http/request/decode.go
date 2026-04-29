// Package request provides utilities for parsing incoming HTTP requests.
// It handles JSON decoding and structural validation of incoming payloads.
package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

// DecodeAndValidateRequest decodes the JSON body of an HTTP request into the provided
// destination struct and validates its fields based on struct tags.
// It returns an ErrInvallidArgument if decoding or validation fails.
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			errs.ErrInvalidArgument,
		)
	}

	var (
		err error
	)

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			errs.ErrInvalidArgument,
		)
	}

	return nil
}
