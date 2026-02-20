package team

import (
	"context"
	"database/sql"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// ListByUser lists all teams where the user is a member.
func (r *Repository) ListByUser(ctx context.Context, userID user.ID) ([]*team.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_by, t.created_at, t.updated_at
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = ?
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var teams []*team.Team
	for rows.Next() {
		var t team.Team
		var description sql.NullString

		err := rows.Scan(&t.ID, &t.Name, &description, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			t.Description = description.String
		}

		teams = append(teams, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
