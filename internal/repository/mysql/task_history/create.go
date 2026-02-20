package task_history

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Create creates a new task history entry.
func (r *Repository) Create(ctx context.Context, h *task.History) (task.HistoryID, error) {
	query := `
		INSERT INTO task_history (task_id, changed_by, field_name, old_value, new_value, changed_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		h.TaskID, h.ChangedBy, h.FieldName, h.OldValue, h.NewValue, h.ChangedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return task.HistoryID(id), nil //nolint:gosec // ID is always positive from MySQL
}
