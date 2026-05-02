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
	fullNameLen := utf8.RuneCountInString(u.FullName)
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLen,
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

// UserPatch represents the data used to partitially update an existing User.
// Only the fields with Set=true are applied during the patch operation.
type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// NewUserPatch creates a new instance of UserPatch.
func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

// Validate checks whether the UserPatch data corresponds to
// the required fields of User entity.
func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			errs.ErrInvalidArgument,
		)
	}

	return nil
}

// ApplyPatch modifies the User entity using the provided UserPatch data.
// It validates both the patch data and the resulting user state before applying changes.
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
