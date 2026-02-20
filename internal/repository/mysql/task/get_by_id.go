package task

import (
	"context"
	"database/sql"
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetByID retrieves a task by ID.
func (r *Repository) GetByID(ctx context.Context, id task.ID) (*task.Task, error) {
	query := `
		SELECT id, team_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	t, err := r.scanTask(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return t, nil
}
