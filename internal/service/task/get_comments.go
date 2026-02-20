package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetTaskComments retrieves task comments.
func (s *Service) GetTaskComments(ctx context.Context, taskID task.ID, requestedBy user.ID) ([]*task.Comment, error) {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if err := s.verifyMembership(ctx, t.TeamID, requestedBy); err != nil {
		return nil, err
	}

	return s.commentRepo.ListByTask(ctx, taskID)
}
