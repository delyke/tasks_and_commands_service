package task

import (
	"context"
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// TaskService defines the interface for task operations.
type TaskService interface {
	CreateTask(ctx context.Context, input CreateTaskInput) (*task.Task, error)
	GetTask(ctx context.Context, id task.ID) (*task.Task, error)
	ListTasks(ctx context.Context, filter *task.Filter) (*ListResult, error)
	UpdateTask(ctx context.Context, input UpdateTaskInput) (*task.Task, error)
	DeleteTask(ctx context.Context, taskID task.ID, deletedBy user.ID) error
	GetTaskHistory(ctx context.Context, taskID task.ID, requestedBy user.ID) ([]*task.History, error)
	AddComment(ctx context.Context, taskID task.ID, userID user.ID, content string) (*task.Comment, error)
	GetTaskComments(ctx context.Context, taskID task.ID, requestedBy user.ID) ([]*task.Comment, error)
	VerifyTaskAccess(ctx context.Context, taskID task.ID, userID user.ID) (*task.Task, error)
}

// Cache interface for task caching.
type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, keys ...string) error
}
