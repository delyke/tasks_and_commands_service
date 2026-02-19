package task

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Update updates a task.
func (r *Repository) Update(ctx context.Context, t *task.Task) error {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?, assignee_id = ?, due_date = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		t.Title, t.Description, t.Status.String(), t.Priority.String(),
		nullableUserID(t.AssigneeID), nullableTime(t.DueDate), t.ID)
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
