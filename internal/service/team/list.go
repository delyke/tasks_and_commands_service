package team

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// ListUserTeams lists all teams where the user is a member.
func (s *Service) ListUserTeams(ctx context.Context, userID user.ID) ([]*team.Team, error) {
	return s.teamRepo.ListByUser(ctx, userID)
}
