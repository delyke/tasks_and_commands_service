package analytics

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// FindOrphanedAssignees finds tasks where assignee is not a team member.
// Validates data integrity using correlated subquery.
func (r *Repository) FindOrphanedAssignees(ctx context.Context) ([]*task.OrphanedTask, error) {
	query := `
		SELECT
			t.id AS task_id,
			t.title AS task_title,
			t.team_id,
			tm.name AS team_name,
			t.assignee_id,
			u.name AS assignee_name
		FROM tasks t
		INNER JOIN teams tm ON t.team_id = tm.id
		INNER JOIN users u ON t.assignee_id = u.id
		WHERE t.assignee_id IS NOT NULL
		AND NOT EXISTS (
			SELECT 1
			FROM team_members mem
			WHERE mem.team_id = t.team_id
			AND mem.user_id = t.assignee_id
		)
		ORDER BY t.team_id, t.id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var orphaned []*task.OrphanedTask
	for rows.Next() {
		var o task.OrphanedTask
		err := rows.Scan(&o.TaskID, &o.TaskTitle, &o.TeamID, &o.TeamName, &o.AssigneeID, &o.AssigneeName)
		if err != nil {
			return nil, err
		}
		orphaned = append(orphaned, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orphaned, nil
}
