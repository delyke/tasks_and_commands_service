package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetMember retrieves a team member.
func (s *Service) GetMember(ctx context.Context, teamID team.ID, userID user.ID) (*team.Member, error) {
	return s.memberRepo.GetByTeamAndUser(ctx, teamID, userID)
}
