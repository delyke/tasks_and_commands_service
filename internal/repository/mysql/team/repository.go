package team

import (
	"database/sql"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

var _ team.Repository = (*Repository)(nil)

// Repository implements team.Repository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewTeamRepository creates a new TeamRepository.
func NewTeamRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
