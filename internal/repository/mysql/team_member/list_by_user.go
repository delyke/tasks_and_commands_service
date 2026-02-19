package team_member

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// ListByUser lists all memberships for a user.
func (r *Repository) ListByUser(ctx context.Context, userID user.ID) ([]*team.Member, error) {
	query := `
		SELECT id, team_id, user_id, role, joined_at
		FROM team_members
		WHERE user_id = ?
		ORDER BY joined_at ASC
	`
	return r.queryMembers(ctx, query, userID)
}
