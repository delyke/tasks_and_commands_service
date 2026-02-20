package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// DeleteTask deletes a task.
func (s *Service) DeleteTask(ctx context.Context, taskID task.ID, deletedBy user.ID) error {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	if err := s.verifyMembership(ctx, t.TeamID, deletedBy); err != nil {
		return err
	}

	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		return err
	}

	s.invalidateTeamCache(ctx, t.TeamID)

	return nil
}
