package handler

import (
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

func toAPIUser(u *user.User) api.User {
	return api.User{
		ID:        int64(u.ID), //nolint:gosec // safe: domain ID is always positive
		Username:  u.Name,
		Email:     u.Email.String(),
		CreatedAt: u.CreatedAt,
	}
}

func toAPITeam(t *team.Team) *api.Team {
	result := &api.Team{
		ID:        int64(t.ID), //nolint:gosec // safe: domain ID is always positive
		Name:      t.Name,
		CreatedBy: int64(t.CreatedBy), //nolint:gosec // safe: domain ID is always positive
		CreatedAt: t.CreatedAt,
	}
	if t.Description != "" {
		result.Description.SetTo(t.Description)
	}
	return result
}

func toAPITask(t *task.Task) *api.Task {
	result := &api.Task{
		ID:        int64(t.ID), //nolint:gosec // safe: domain ID is always positive
		Title:     t.Title,
		Status:    api.TaskStatus(t.Status),
		Priority:  api.TaskPriority(t.Priority),
		TeamID:    int64(t.TeamID),    //nolint:gosec // safe: domain ID is always positive
		CreatedBy: int64(t.CreatedBy), //nolint:gosec // safe: domain ID is always positive
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if t.Description != "" {
		result.Description.SetTo(t.Description)
	}
	if t.AssigneeID != nil {
		result.AssigneeID.SetTo(int64(*t.AssigneeID)) //nolint:gosec // safe: domain ID is always positive
	}
	if t.DueDate != nil {
		result.DueDate.SetTo(*t.DueDate)
	}
	return result
}
