package task_comment

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Update updates a comment.
func (r *Repository) Update(ctx context.Context, c *task.Comment) error {
	query := `
		UPDATE task_comments
		SET content = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, c.Content, c.ID)
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
