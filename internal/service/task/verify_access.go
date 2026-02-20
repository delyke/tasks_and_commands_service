package task

import (
	"context"
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// VerifyTaskAccess verifies that a user has access to a task.
func (s *Service) VerifyTaskAccess(ctx context.Context, taskID task.ID, userID user.ID) (*task.Task, error) {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if err := s.verifyMembership(ctx, t.TeamID, userID); err != nil {
		return nil, err
	}

	return t, nil
}
