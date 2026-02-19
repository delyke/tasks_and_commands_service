package team_member

import (
	"context"
	"database/sql"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

var _ team.MemberRepository = (*Repository)(nil)

// Repository implements team.MemberRepository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewTeamMemberRepository creates a new Repository.
func NewTeamMemberRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) queryMembers(ctx context.Context, query string, arg any) ([]*team.Member, error) {
	rows, err := r.db.QueryContext(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var members []*team.Member
	for rows.Next() {
		var m team.Member
		var role string

		err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &role, &m.JoinedAt)
		if err != nil {
			return nil, err
		}

		m.Role = team.Role(role)
		members = append(members, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
