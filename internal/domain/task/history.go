package task

import (
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// HistoryID represents a task history entry identifier.
type HistoryID uint64

// History represents a task change history entry.
type History struct {
	ID        HistoryID
	TaskID    ID
	ChangedBy user.ID
	FieldName string
	OldValue  *string
	NewValue  *string
	ChangedAt time.Time
}

// NewHistory creates a new History entry.
func NewHistory(taskID ID, changedBy user.ID, fieldName string, oldValue, newValue *string) *History {
	return &History{
		TaskID:    taskID,
		ChangedBy: changedBy,
		FieldName: fieldName,
		OldValue:  oldValue,
		NewValue:  newValue,
		ChangedAt: time.Now(),
	}
}
