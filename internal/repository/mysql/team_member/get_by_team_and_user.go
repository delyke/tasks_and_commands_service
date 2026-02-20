package team_member

import (
	"context"
	"database/sql"
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetByTeamAndUser retrieves a member by team and user ID.
func (r *Repository) GetByTeamAndUser(ctx context.Context, teamID team.ID, userID user.ID) (*team.Member, error) {
	query := `
		SELECT id, team_id, user_id, role, joined_at
		FROM team_members
		WHERE team_id = ? AND user_id = ?
	`

	var m team.Member
	var role string

	err := r.db.QueryRowContext(ctx, query, teamID, userID).Scan(
		&m.ID, &m.TeamID, &m.UserID, &role, &m.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	m.Role = team.Role(role)
	return &m, nil
}
