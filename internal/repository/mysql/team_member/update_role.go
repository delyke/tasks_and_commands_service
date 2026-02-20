package team_member

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// UpdateRole updates a member's role.
func (r *Repository) UpdateRole(ctx context.Context, id team.MemberID, role team.Role) error {
	query := `UPDATE team_members SET role = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, role.String(), id)
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
