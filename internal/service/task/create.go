package task

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"time"
)

// CreateTaskInput represents task creation input.
type CreateTaskInput struct {
	TeamID      team.ID
	Title       string
	Description string
	Priority    string
	AssigneeID  *user.ID
	DueDate     *time.Time
	CreatedBy   user.ID
}

// CreateTask creates a new task.
func (s *Service) CreateTask(ctx context.Context, input CreateTaskInput) (*task.Task, error) {
	if err := s.verifyMembership(ctx, input.TeamID, input.CreatedBy); err != nil {
		return nil, err
	}

	if input.AssigneeID != nil {
		if err := s.verifyAssignee(ctx, input.TeamID, *input.AssigneeID); err != nil {
			return nil, err
		}
	}

	priority, err := task.NewPriority(input.Priority)
	if err != nil {
		priority = task.PriorityMedium
	}

	t := task.New(input.TeamID, input.Title, input.Description, priority, input.CreatedBy, input.AssigneeID, input.DueDate)

	taskID, err := s.taskRepo.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	t = t.WithID(taskID)

	s.invalidateTeamCache(ctx, input.TeamID)

	return t, nil
}
