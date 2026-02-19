package task_comment

import (
	"context"
	"database/sql"
	"errors"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetByID retrieves a comment by ID.
func (r *Repository) GetByID(ctx context.Context, id task.CommentID) (*task.Comment, error) {
	query := `
		SELECT id, task_id, user_id, content, created_at, updated_at
		FROM task_comments
		WHERE id = ?
	`

	var c task.Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &c, nil
}
