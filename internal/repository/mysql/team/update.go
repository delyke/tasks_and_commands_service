package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// Update updates a team.
func (r *Repository) Update(ctx context.Context, t *team.Team) error {
	query := `
		UPDATE teams
		SET name = ?, description = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, t.Name, t.Description, t.ID)
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
