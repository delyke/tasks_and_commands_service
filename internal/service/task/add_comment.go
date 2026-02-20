package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// AddComment adds a comment to a task.
func (s *Service) AddComment(ctx context.Context, taskID task.ID, userID user.ID, content string) (*task.Comment, error) {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if err := s.verifyMembership(ctx, t.TeamID, userID); err != nil {
		return nil, err
	}

	comment := task.NewComment(taskID, userID, content)
	commentID, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment.WithID(commentID), nil
}
