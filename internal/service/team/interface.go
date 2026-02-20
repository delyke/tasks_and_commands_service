package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// TeamService defines the interface for team operations.
type TeamService interface {
	CreateTeam(ctx context.Context, input CreateTeamInput) (*team.Team, error)
	ListUserTeams(ctx context.Context, userID user.ID) ([]*team.Team, error)
	InviteMember(ctx context.Context, input InviteInput) error
	GetMember(ctx context.Context, teamID team.ID, userID user.ID) (*team.Member, error)
	IsMember(ctx context.Context, teamID team.ID, userID user.ID) (bool, error)
}
