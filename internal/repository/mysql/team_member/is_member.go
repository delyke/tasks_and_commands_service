package team_member

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// IsMember checks if a user is a member of a team.
func (r *Repository) IsMember(ctx context.Context, teamID team.ID, userID user.ID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = ? AND user_id = ?)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, teamID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
