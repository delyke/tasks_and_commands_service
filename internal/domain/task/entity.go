package task

import (
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// ID represents a task identifier.
type ID uint64

// Task represents a domain task entity.
type Task struct {
	ID          ID
	TeamID      team.ID
	Title       string
	Description string
	Status      Status
	Priority    Priority
	AssigneeID  *user.ID
	CreatedBy   user.ID
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// New creates a new Task entity.
func New(
	teamID team.ID,
	title, description string,
	priority Priority,
	createdBy user.ID,
	assigneeID *user.ID,
	dueDate *time.Time,
) *Task {
	now := time.Now()
	return &Task{
		TeamID:      teamID,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		Priority:    priority,
		AssigneeID:  assigneeID,
		CreatedBy:   createdBy,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// WithID returns a copy of the task with the given ID.
func (t *Task) WithID(id ID) *Task {
	t.ID = id
	return t
}

// UpdateStatus updates the task status and records the change time.
func (t *Task) UpdateStatus(status Status) {
	t.Status = status
	t.UpdatedAt = time.Now()
}

// UpdateAssignee updates the task assignee.
func (t *Task) UpdateAssignee(assigneeID *user.ID) {
	t.AssigneeID = assigneeID
	t.UpdatedAt = time.Now()
}
