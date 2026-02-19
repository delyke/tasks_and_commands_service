package team

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// ListTeamMembers lists all members of a team.
func (s *Service) ListTeamMembers(ctx context.Context, teamID team.ID) ([]*team.Member, error) {
	return s.memberRepo.ListByTeam(ctx, teamID)
}
