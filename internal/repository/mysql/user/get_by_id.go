package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetByID retrieves a user by ID.
func (r *Repository) GetByID(ctx context.Context, id user.ID) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var u user.User
	var email string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	u.Email = user.Email(email)
	return &u, nil
}
