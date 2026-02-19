package team

import (
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// MemberID represents a team member identifier.
type MemberID uint64

// Member represents a team membership entity.
type Member struct {
	ID       MemberID
	TeamID   ID
	UserID   user.ID
	Role     Role
	JoinedAt time.Time
}

// NewMember creates a new Member entity.
func NewMember(teamID ID, userID user.ID, role Role) *Member {
	return &Member{
		TeamID:   teamID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

// WithID returns a copy of the member with the given ID.
func (m *Member) WithID(id MemberID) *Member {
	m.ID = id
	return m
}
