package task_history

import (
	"database/sql"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

var _ task.HistoryRepository = (*Repository)(nil)

// Repository implements task.HistoryRepository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewTaskHistoryRepository creates a new Repository.
func NewTaskHistoryRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
