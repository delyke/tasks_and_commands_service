package team

import (
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// ID represents a team identifier.
type ID uint64

// Team represents a domain team entity.
type Team struct {
	ID          ID
	Name        string
	Description string
	CreatedBy   user.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New creates a new Team entity.
func New(name, description string, createdBy user.ID) *Team {
	now := time.Now()
	return &Team{
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// WithID returns a copy of the team with the given ID.
func (t *Team) WithID(id ID) *Team {
	t.ID = id
	return t
}
