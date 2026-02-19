package task

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Create creates a new task.
func (r *Repository) Create(ctx context.Context, t *task.Task) (task.ID, error) {
	query := `
		INSERT INTO tasks (team_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		t.TeamID, t.Title, t.Description, t.Status.String(), t.Priority.String(),
		nullableUserID(t.AssigneeID), t.CreatedBy, nullableTime(t.DueDate),
		t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return task.ID(id), nil //nolint:gosec // ID is always positive from MySQL
}
