package task

import (
	"database/sql"
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Ensure interface compliance
var _ task.Repository = (*Repository)(nil)

// Repository implements task.Repository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// scanner is an interface for scanning database rows.
type scanner interface {
	Scan(dest ...any) error
}

func (r *Repository) scanTask(row scanner) (*task.Task, error) {
	var t task.Task
	var description sql.NullString
	var status, priority string
	var assigneeID sql.NullInt64
	var dueDate sql.NullTime

	err := row.Scan(
		&t.ID, &t.TeamID, &t.Title, &description, &status, &priority,
		&assigneeID, &t.CreatedBy, &dueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		t.Description = description.String
	}
	t.Status = task.Status(status)
	t.Priority = task.Priority(priority)
	if assigneeID.Valid {
		id := user.ID(assigneeID.Int64) //nolint:gosec // ID is always positive
		t.AssigneeID = &id
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}

	return &t, nil
}

func (r *Repository) scanTasks(rows *sql.Rows) ([]*task.Task, error) {
	defer func() { _ = rows.Close() }() //nolint:errcheck,gosec // rows.Close() error is non-critical

	var tasks []*task.Task
	for rows.Next() {
		t, err := r.scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func nullableUserID(id *user.ID) sql.NullInt64 {
	if id == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*id), Valid: true} //nolint:gosec // ID is always positive
}

func nullableTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
