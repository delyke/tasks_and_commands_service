package task

import "errors"

var ErrInvalidStatus = errors.New("invalid task status")

// Status represents a task status.
type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// NewStatus creates and validates a Status value object.
func NewStatus(value string) (Status, error) {
	switch Status(value) {
	case StatusTodo, StatusInProgress, StatusDone:
		return Status(value), nil
	default:
		return "", ErrInvalidStatus
	}
}

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}
