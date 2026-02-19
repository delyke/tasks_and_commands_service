package analytics

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetTopCreatorsByTeam returns top task creators per team.
func (s *Service) GetTopCreatorsByTeam(ctx context.Context, limit int) ([]*task.TopCreator, error) {
	if limit <= 0 {
		limit = 3
	}
	return s.analyticsRepo.GetTopCreatorsByTeam(ctx, limit)
}
