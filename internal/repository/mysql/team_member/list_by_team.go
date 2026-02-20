package team_member

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// ListByTeam lists all members of a team.
func (r *Repository) ListByTeam(ctx context.Context, teamID team.ID) ([]*team.Member, error) {
	query := `
		SELECT id, team_id, user_id, role, joined_at
		FROM team_members
		WHERE team_id = ?
		ORDER BY joined_at ASC
	`
	return r.queryMembers(ctx, query, teamID)
}
