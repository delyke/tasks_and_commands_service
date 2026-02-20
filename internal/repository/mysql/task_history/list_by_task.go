package task_history

import (
	"context"
	"database/sql"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// ListByTask lists all history entries for a task.
func (r *Repository) ListByTask(ctx context.Context, taskID task.ID) ([]*task.History, error) {
	query := `
		SELECT id, task_id, changed_by, field_name, old_value, new_value, changed_at
		FROM task_history
		WHERE task_id = ?
		ORDER BY changed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var history []*task.History
	for rows.Next() {
		var h task.History
		var oldValue, newValue sql.NullString

		err := rows.Scan(&h.ID, &h.TaskID, &h.ChangedBy, &h.FieldName, &oldValue, &newValue, &h.ChangedAt)
		if err != nil {
			return nil, err
		}

		if oldValue.Valid {
			h.OldValue = &oldValue.String
		}
		if newValue.Valid {
			h.NewValue = &newValue.String
		}

		history = append(history, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
