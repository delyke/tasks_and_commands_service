package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Repository defines the interface for team persistence.
type Repository interface {
	Create(ctx context.Context, team *Team) (ID, error)
	GetByID(ctx context.Context, id ID) (*Team, error)
	Update(ctx context.Context, team *Team) error
	Delete(ctx context.Context, id ID) error
	ListByUser(ctx context.Context, userID user.ID) ([]*Team, error)
}

// MemberRepository defines the interface for team member persistence.
type MemberRepository interface {
	Add(ctx context.Context, member *Member) (MemberID, error)
	GetByTeamAndUser(ctx context.Context, teamID ID, userID user.ID) (*Member, error)
	ListByTeam(ctx context.Context, teamID ID) ([]*Member, error)
	ListByUser(ctx context.Context, userID user.ID) ([]*Member, error)
	UpdateRole(ctx context.Context, id MemberID, role Role) error
	Remove(ctx context.Context, teamID ID, userID user.ID) error
	IsMember(ctx context.Context, teamID ID, userID user.ID) (bool, error)
}
