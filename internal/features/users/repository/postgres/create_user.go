package postgres

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
)

// CreateUser executes the SQL query to insert a new user into the database.
// It applies a configured operation timeout, maps the resulting database row
// back into a domain enitity.
func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		user.FullName,
		user.PhoneNumber,
	)

	return userDomain, nil
}
