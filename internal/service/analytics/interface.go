package analytics

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// AnalyticsService defines the interface for analytics operations.
type AnalyticsService interface {
	GetTeamStats(ctx context.Context) ([]*task.TeamStats, error)
	GetTopCreatorsByTeam(ctx context.Context, limit int) ([]*task.TopCreator, error)
	FindOrphanedAssignees(ctx context.Context) ([]*task.OrphanedTask, error)
}
