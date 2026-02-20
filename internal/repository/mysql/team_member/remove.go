package team_member

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Remove removes a member from a team.
func (r *Repository) Remove(ctx context.Context, teamID team.ID, userID user.ID) error {
	query := `DELETE FROM team_members WHERE team_id = ? AND user_id = ?`

	result, err := r.db.ExecContext(ctx, query, teamID, userID)
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
