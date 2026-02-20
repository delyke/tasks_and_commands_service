package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// GetTask retrieves a task by ID.
func (s *Service) GetTask(ctx context.Context, id task.ID) (*task.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}
