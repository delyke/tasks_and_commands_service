package task_comment

import (
	"database/sql"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Ensure interface compliance
var _ task.CommentRepository = (*Repository)(nil)

// Repository implements task.CommentRepository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewTaskCommentRepository creates a new Repository.
func NewTaskCommentRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
