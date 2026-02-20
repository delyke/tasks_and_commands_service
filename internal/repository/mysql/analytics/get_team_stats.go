package analytics

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetTeamStats returns team statistics with JOIN 3+ tables and aggregation.
// Query: Team name, member count, done tasks in last 7 days.
func (r *Repository) GetTeamStats(ctx context.Context) ([]*task.TeamStats, error) {
	query := `
		SELECT
			t.id AS team_id,
			t.name AS team_name,
			COUNT(DISTINCT tm.user_id) AS member_count,
			COUNT(DISTINCT CASE
				WHEN tk.status = 'done'
				AND tk.updated_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
				THEN tk.id
			END) AS done_tasks_last_7d
		FROM teams t
		LEFT JOIN team_members tm ON t.id = tm.team_id
		LEFT JOIN tasks tk ON t.id = tk.team_id
		GROUP BY t.id, t.name
		ORDER BY t.name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var stats []*task.TeamStats
	for rows.Next() {
		var s task.TeamStats
		err := rows.Scan(&s.TeamID, &s.TeamName, &s.MemberCount, &s.DoneTasksLast7d)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
