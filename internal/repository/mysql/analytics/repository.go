package analytics

import (
	"database/sql"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Ensure interface compliance
var _ task.AnalyticsRepository = (*Repository)(nil)

// Repository implements task.AnalyticsRepository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewAnalyticsRepository creates a new AnalyticsRepository.
func NewAnalyticsRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
