package task

import "errors"

var ErrInvalidPriority = errors.New("invalid task priority")

// Priority represents a task priority.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// NewPriority creates and validates a Priority value object.
func NewPriority(value string) (Priority, error) {
	switch Priority(value) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return Priority(value), nil
	default:
		return "", ErrInvalidPriority
	}
}

// String returns the string representation of the priority.
func (p Priority) String() string {
	return string(p)
}
