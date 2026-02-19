package user

import (
	"database/sql"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

var _ user.Repository = (*Repository)(nil)

// Repository implements user.Repository using MySQL.
type Repository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
