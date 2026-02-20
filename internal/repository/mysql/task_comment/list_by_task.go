package task_comment

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// ListByTask lists all comments for a task.
func (r *Repository) ListByTask(ctx context.Context, taskID task.ID) ([]*task.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM task_comments
		WHERE task_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var comments []*task.Comment
	for rows.Next() {
		var c task.Comment
		err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
