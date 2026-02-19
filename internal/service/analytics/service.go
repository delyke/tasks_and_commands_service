package analytics

import (
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

// Service handles analytics operations.
type Service struct {
	analyticsRepo task.AnalyticsRepository
}

// NewAnalyticsService creates a new AnalyticsService.
func NewAnalyticsService(analyticsRepo task.AnalyticsRepository) *Service {
	return &Service{
		analyticsRepo: analyticsRepo,
	}
}
