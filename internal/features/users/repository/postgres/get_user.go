package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id=`%d`: %w",
				id,
				errs.ErrNotFound,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
