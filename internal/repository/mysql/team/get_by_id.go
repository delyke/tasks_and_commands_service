package team

import (
	"context"
	"database/sql"
	"errors"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// GetByID retrieves a team by ID.
func (r *Repository) GetByID(ctx context.Context, id team.ID) (*team.Team, error) {
	query := `
		SELECT id, name, description, created_by, created_at, updated_at
		FROM teams
		WHERE id = ?
	`

	var t team.Team
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Name, &description, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if description.Valid {
		t.Description = description.String
	}

	return &t, nil
}
