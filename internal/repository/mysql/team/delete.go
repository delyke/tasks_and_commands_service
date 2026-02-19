package team

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// Delete deletes a team.
func (r *Repository) Delete(ctx context.Context, id team.ID) error {
	query := `DELETE FROM teams WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
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
