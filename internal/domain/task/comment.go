package task

import (
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// CommentID represents a task comment identifier.
type CommentID uint64

// Comment represents a task comment entity.
type Comment struct {
	ID        CommentID
	TaskID    ID
	UserID    user.ID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewComment creates a new Comment entity.
func NewComment(taskID ID, userID user.ID, content string) *Comment {
	now := time.Now()
	return &Comment{
		TaskID:    taskID,
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithID returns a copy of the comment with the given ID.
func (c *Comment) WithID(id CommentID) *Comment {
	c.ID = id
	return c
}

// UpdateContent updates the comment content.
func (c *Comment) UpdateContent(content string) {
	c.Content = content
	c.UpdatedAt = time.Now()
}
