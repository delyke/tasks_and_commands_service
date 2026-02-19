package task_comment

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Create creates a new comment.
func (r *Repository) Create(ctx context.Context, c *task.Comment) (task.CommentID, error) {
	query := `
		INSERT INTO task_comments (task_id, user_id, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		c.TaskID, c.UserID, c.Content, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return task.CommentID(id), nil //nolint:gosec // ID is always positive from MySQL
}
