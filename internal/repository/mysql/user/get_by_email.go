package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetByEmail retrieves a user by email.
func (r *Repository) GetByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var u user.User
	var emailStr string

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&u.ID, &emailStr, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	u.Email = user.Email(emailStr)
	return &u, nil
}
