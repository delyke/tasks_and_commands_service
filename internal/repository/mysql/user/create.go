package user

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Create creates a new user.
func (r *Repository) Create(ctx context.Context, u *user.User) (user.ID, error) {
	query := `
		INSERT INTO users (email, password_hash, name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		u.Email.String(), u.PasswordHash, u.Name, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return user.ID(id), nil //nolint:gosec // ID is always positive from MySQL
}
