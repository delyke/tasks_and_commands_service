package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

// ListResult represents a paginated task list.
type ListResult struct {
	Tasks []*task.Task
	Total int64
}

// ListTasks lists tasks with filtering and pagination.
func (s *Service) ListTasks(ctx context.Context, filter *task.Filter) (*ListResult, error) {
	if s.isCacheable(filter) {
		cacheKey := s.teamCacheKey(team.ID(*filter.TeamID))
		var cached ListResult
		if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
			return &cached, nil
		}
	}

	tasks, total, err := s.taskRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := &ListResult{
		Tasks: tasks,
		Total: total,
	}

	if s.isCacheable(filter) {
		cacheKey := s.teamCacheKey(team.ID(*filter.TeamID))
		_ = s.cache.Set(ctx, cacheKey, result, tasksCacheTTL) //nolint:errcheck,gosec // cache errors are non-critical
	}

	return result, nil
}
