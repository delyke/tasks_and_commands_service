package task

import (
	"context"
	"strings"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// List lists tasks with filtering and pagination.
func (r *Repository) List(ctx context.Context, filter *task.Filter) ([]*task.Task, int64, error) {
	var conditions []string
	var args []any

	if filter.TeamID != nil {
		conditions = append(conditions, "team_id = ?")
		args = append(args, *filter.TeamID)
	}
	if filter.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status.String())
	}
	if filter.AssigneeID != nil {
		conditions = append(conditions, "assignee_id = ?")
		args = append(args, *filter.AssigneeID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) FROM tasks " + whereClause //nolint:gosec // safe concatenation with validated conditions
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	query := "SELECT id, team_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at FROM tasks " + //nolint:gosec // safe concatenation with validated conditions
		whereClause + " ORDER BY created_at DESC LIMIT ? OFFSET ?"

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	tasks, err := r.scanTasks(rows)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
