package user

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Update updates a user.
func (r *Repository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET email = ?, password_hash = ?, name = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		u.Email.String(), u.PasswordHash, u.Name, u.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
