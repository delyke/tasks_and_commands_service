package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// ListByTeam lists tasks for a specific team.
func (r *Repository) ListByTeam(ctx context.Context, teamID team.ID, limit, offset int) ([]*task.Task, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM tasks WHERE team_id = ?`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, teamID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	query := `
		SELECT id, team_id, title, description, status, priority, assignee_id, created_by, due_date, created_at, updated_at
		FROM tasks
		WHERE team_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, teamID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	tasks, err := r.scanTasks(rows)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
