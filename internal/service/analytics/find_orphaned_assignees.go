package analytics

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// FindOrphanedAssignees finds tasks with orphaned assignees.
func (s *Service) FindOrphanedAssignees(ctx context.Context) ([]*task.OrphanedTask, error) {
	return s.analyticsRepo.FindOrphanedAssignees(ctx)
}
