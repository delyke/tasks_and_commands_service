package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// CreateTeamInput represents team creation input.
type CreateTeamInput struct {
	Name        string
	Description string
	CreatedBy   user.ID
}

// CreateTeam creates a new team and adds the creator as owner.
func (s *Service) CreateTeam(ctx context.Context, input CreateTeamInput) (*team.Team, error) {
	t := team.New(input.Name, input.Description, input.CreatedBy)

	teamID, err := s.teamRepo.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	t = t.WithID(teamID)

	// Add creator as owner
	member := team.NewMember(teamID, input.CreatedBy, team.RoleOwner)
	_, err = s.memberRepo.Add(ctx, member)
	if err != nil {
		return nil, err
	}

	return t, nil
}
