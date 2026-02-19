package task

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetTaskHistory retrieves task change history.
func (s *Service) GetTaskHistory(ctx context.Context, taskID task.ID, requestedBy user.ID) ([]*task.History, error) {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if err := s.verifyMembership(ctx, t.TeamID, requestedBy); err != nil {
		return nil, err
	}

	return s.historyRepo.ListByTask(ctx, taskID)
}
