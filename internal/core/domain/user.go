package domain

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// User represents the core business entity of a user in the system.
// It contains all the essential data and business logic tied to a user.
type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

// NewUser reconstitutes an existing User entity from storage
// with a known ID and Version.
func NewUser(id, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

var phoneRegex = regexp.MustCompile(`^\+[0-9]+$`)

// Validate checks whether the business rules for the User entity are met.
// It returns error if the data contradicts the rules
// such as length or symbols type of full_name and phone_number.
func (u *User) Validate() error {
	fullNameLength := utf8.RuneCountInString(u.FullName)
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			errs.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := utf8.RuneCountInString(*u.PhoneNumber)
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				errs.ErrInvalidArgument,
			)
		}

		if !phoneRegex.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				errs.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			errs.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u
	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp
	return nil
}
