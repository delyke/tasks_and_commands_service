package analytics

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetTeamStats returns team statistics.
func (s *Service) GetTeamStats(ctx context.Context) ([]*task.TeamStats, error) {
	return s.analyticsRepo.GetTeamStats(ctx)
}
