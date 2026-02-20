package task

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

const (
	tasksCacheKeyPrefix = "tasks:team:"
	tasksCacheTTL       = 5 * time.Minute
)

// Service handles task operations.
type Service struct {
	taskRepo    task.Repository
	historyRepo task.HistoryRepository
	commentRepo task.CommentRepository
	memberRepo  team.MemberRepository
	cache       Cache
}

// Ensure Service implements TaskService.
var _ TaskService = (*Service)(nil)

// NewTaskService creates a new TaskService.
func NewTaskService(
	taskRepo task.Repository,
	historyRepo task.HistoryRepository,
	commentRepo task.CommentRepository,
	memberRepo team.MemberRepository,
	cache Cache,
) *Service {
	return &Service{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
		commentRepo: commentRepo,
		memberRepo:  memberRepo,
		cache:       cache,
	}
}

type historyChange struct {
	field    string
	oldValue string
	newValue string
}

func (s *Service) verifyMembership(ctx context.Context, teamID team.ID, userID user.ID) error {
	isMember, err := s.memberRepo.IsMember(ctx, teamID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.ErrForbidden
	}
	return nil
}

func (s *Service) verifyAssignee(ctx context.Context, teamID team.ID, assigneeID user.ID) error {
	isMember, err := s.memberRepo.IsMember(ctx, teamID, assigneeID)
	if err != nil {
		return err
	}
	if !isMember {
		return fmt.Errorf("assignee is not a team member: %w", domain.ErrInvalidInput)
	}
	return nil
}

func (s *Service) isCacheable(filter *task.Filter) bool {
	return filter.TeamID != nil && filter.Status == nil && filter.AssigneeID == nil && filter.Offset == 0
}

func (s *Service) applyTaskChanges(ctx context.Context, t *task.Task, input UpdateTaskInput) ([]historyChange, error) {
	var changes []historyChange

	changes = s.applyStringChange(changes, "title", t.Title, input.Title, func(v string) { t.Title = v })
	changes = s.applyStringChange(changes, "description", t.Description, input.Description, func(v string) { t.Description = v })

	if input.Status != nil {
		if newStatus, err := task.NewStatus(*input.Status); err == nil && newStatus != t.Status {
			changes = append(changes, historyChange{field: "status", oldValue: t.Status.String(), newValue: newStatus.String()})
			t.Status = newStatus
		}
	}

	if input.Priority != nil {
		if newPriority, err := task.NewPriority(*input.Priority); err == nil && newPriority != t.Priority {
			changes = append(changes, historyChange{field: "priority", oldValue: t.Priority.String(), newValue: newPriority.String()})
			t.Priority = newPriority
		}
	}

	if input.AssigneeID != nil {
		if err := s.verifyAssignee(ctx, t.TeamID, *input.AssigneeID); err != nil {
			return nil, err
		}
		changes = s.applyAssigneeChange(changes, t, input.AssigneeID)
	}

	return changes, nil
}

func (s *Service) applyStringChange(changes []historyChange, field, oldVal string, newVal *string, apply func(string)) []historyChange {
	if newVal != nil && *newVal != oldVal {
		changes = append(changes, historyChange{field: field, oldValue: oldVal, newValue: *newVal})
		apply(*newVal)
	}
	return changes
}

func (s *Service) applyAssigneeChange(changes []historyChange, t *task.Task, newAssigneeID *user.ID) []historyChange {
	oldAssignee := ""
	if t.AssigneeID != nil {
		oldAssignee = strconv.FormatUint(uint64(*t.AssigneeID), 10)
	}
	newAssignee := strconv.FormatUint(uint64(*newAssigneeID), 10)
	if oldAssignee != newAssignee {
		changes = append(changes, historyChange{field: "assignee_id", oldValue: oldAssignee, newValue: newAssignee})
		t.AssigneeID = newAssigneeID
	}
	return changes
}

func (s *Service) recordHistory(ctx context.Context, taskID task.ID, changedBy user.ID, changes []historyChange) {
	for _, change := range changes {
		h := task.NewHistory(taskID, changedBy, change.field, &change.oldValue, &change.newValue)
		_, _ = s.historyRepo.Create(ctx, h) //nolint:errcheck,gosec // history errors are non-critical
	}
}

func (s *Service) teamCacheKey(teamID team.ID) string {
	return fmt.Sprintf("%s%d", tasksCacheKeyPrefix, teamID)
}

func (s *Service) invalidateTeamCache(ctx context.Context, teamID team.ID) {
	_ = s.cache.Delete(ctx, s.teamCacheKey(teamID)) //nolint:errcheck,gosec // cache errors are non-critical
}
