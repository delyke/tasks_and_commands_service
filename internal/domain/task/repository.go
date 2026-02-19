package task

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Repository defines the interface for task persistence.
type Repository interface {
	Create(ctx context.Context, task *Task) (ID, error)
	GetByID(ctx context.Context, id ID) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context, filter *Filter) ([]*Task, int64, error)
	ListByTeam(ctx context.Context, teamID team.ID, limit, offset int) ([]*Task, int64, error)
}

// HistoryRepository defines the interface for task history persistence.
type HistoryRepository interface {
	Create(ctx context.Context, history *History) (HistoryID, error)
	ListByTask(ctx context.Context, taskID ID) ([]*History, error)
}

// CommentRepository defines the interface for task comment persistence.
type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) (CommentID, error)
	GetByID(ctx context.Context, id CommentID) (*Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id CommentID) error
	ListByTask(ctx context.Context, taskID ID) ([]*Comment, error)
}

// AnalyticsRepository defines complex query methods.
type AnalyticsRepository interface {
	// GetTeamStats returns team statistics (JOIN 3+ tables + aggregation)
	// Team name, member count, done tasks in last 7 days
	GetTeamStats(ctx context.Context) ([]*TeamStats, error)

	// GetTopCreatorsByTeam returns top-3 task creators per team (window function)
	GetTopCreatorsByTeam(ctx context.Context, limit int) ([]*TopCreator, error)

	// FindOrphanedAssignees finds tasks where assignee is not a team member
	FindOrphanedAssignees(ctx context.Context) ([]*OrphanedTask, error)
}

// TeamStats represents team statistics.
type TeamStats struct {
	TeamID          team.ID
	TeamName        string
	MemberCount     int64
	DoneTasksLast7d int64
}

// TopCreator represents a top task creator.
type TopCreator struct {
	TeamID    team.ID
	TeamName  string
	UserID    user.ID
	UserName  string
	TaskCount int64
	Rank      int
}

// OrphanedTask represents a task with orphaned assignee.
type OrphanedTask struct {
	TaskID       ID
	TaskTitle    string
	TeamID       team.ID
	TeamName     string
	AssigneeID   user.ID
	AssigneeName string
}
