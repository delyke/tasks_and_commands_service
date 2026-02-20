package team_member

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// Add adds a new team member.
func (r *Repository) Add(ctx context.Context, m *team.Member) (team.MemberID, error) {
	query := `
		INSERT INTO team_members (team_id, user_id, role, joined_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		m.TeamID, m.UserID, m.Role.String(), m.JoinedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return team.MemberID(id), nil //nolint:gosec // ID is always positive from MySQL
}
