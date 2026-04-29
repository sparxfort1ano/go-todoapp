package postgres

import "github.com/sparxfort1ano/go-todoapp/internal/core/domain"

// UserModel represents the database schema for a user (DAO).
type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

// userDomainsFromModels creates a new slice of domain.User
// with the given user model slice.
func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return userDomains
}
