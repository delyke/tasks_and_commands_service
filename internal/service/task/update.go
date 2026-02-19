package task

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"time"
)

// UpdateTaskInput represents task update input.
type UpdateTaskInput struct {
	TaskID      task.ID
	Title       *string
	Description *string
	Status      *string
	Priority    *string
	AssigneeID  *user.ID
	DueDate     *time.Time
	UpdatedBy   user.ID
}

// UpdateTask updates a task.
func (s *Service) UpdateTask(ctx context.Context, input UpdateTaskInput) (*task.Task, error) {
	t, err := s.taskRepo.GetByID(ctx, input.TaskID)
	if err != nil {
		return nil, err
	}

	if err := s.verifyMembership(ctx, t.TeamID, input.UpdatedBy); err != nil {
		return nil, err
	}

	changes, err := s.applyTaskChanges(ctx, t, input)
	if err != nil {
		return nil, err
	}

	t.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, t); err != nil {
		return nil, err
	}

	s.recordHistory(ctx, t.ID, input.UpdatedBy, changes)
	s.invalidateTeamCache(ctx, t.TeamID)

	return t, nil
}
