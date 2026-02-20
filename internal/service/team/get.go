package team

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// GetTeam retrieves a team by ID.
func (s *Service) GetTeam(ctx context.Context, id team.ID) (*team.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}
