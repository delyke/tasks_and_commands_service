package team

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// IsMember checks if a user is a member of a team.
func (s *Service) IsMember(ctx context.Context, teamID team.ID, userID user.ID) (bool, error) {
	return s.memberRepo.IsMember(ctx, teamID, userID)
}
