package analytics

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetTopCreatorsByTeam returns top-N task creators per team using window functions.
// Uses ROW_NUMBER() window function for ranking.
func (r *Repository) GetTopCreatorsByTeam(ctx context.Context, limit int) ([]*task.TopCreator, error) {
	query := `
		WITH task_counts AS (
			SELECT
				t.team_id,
				tm.name AS team_name,
				t.created_by AS user_id,
				u.name AS user_name,
				COUNT(*) AS task_count,
				ROW_NUMBER() OVER (
					PARTITION BY t.team_id
					ORDER BY COUNT(*) DESC
				) AS rank_num
			FROM tasks t
			INNER JOIN teams tm ON t.team_id = tm.id
			INNER JOIN users u ON t.created_by = u.id
			WHERE t.created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)
			GROUP BY t.team_id, tm.name, t.created_by, u.name
		)
		SELECT team_id, team_name, user_id, user_name, task_count, rank_num
		FROM task_counts
		WHERE rank_num <= ?
		ORDER BY team_id, rank_num
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var creators []*task.TopCreator
	for rows.Next() {
		var c task.TopCreator
		err := rows.Scan(&c.TeamID, &c.TeamName, &c.UserID, &c.UserName, &c.TaskCount, &c.Rank)
		if err != nil {
			return nil, err
		}
		creators = append(creators, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return creators, nil
}
