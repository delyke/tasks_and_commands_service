package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// Create creates a new team.
func (r *Repository) Create(ctx context.Context, t *team.Team) (team.ID, error) {
	query := `
		INSERT INTO teams (name, description, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		t.Name, t.Description, t.CreatedBy, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return team.ID(id), nil //nolint:gosec // ID is always positive from MySQL
}
