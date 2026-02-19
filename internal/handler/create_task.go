package handler

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	"github.com/delyke/tasks_and_commands_service/internal/service/task"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// CreateTask implements api.Handler.
func (h *Handler) CreateTask(ctx context.Context, req *api.CreateTaskReq) (api.CreateTaskRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	var assigneeID *user.ID
	if req.AssigneeID.IsSet() {
		id := user.ID(req.AssigneeID.Value) //nolint:gosec // safe: API ID is validated
		assigneeID = &id
	}

	priority := "medium"
	if req.Priority.IsSet() {
		priority = string(req.Priority.Value)
	}

	input := task.CreateTaskInput{
		TeamID:      team.ID(req.TeamID), //nolint:gosec // safe: API ID is validated
		Title:       req.Title,
		Description: req.Description.Or(""),
		Priority:    priority,
		AssigneeID:  assigneeID,
		CreatedBy:   user.ID(userID), //nolint:gosec // safe: API ID is validated
	}

	if req.DueDate.IsSet() {
		t := req.DueDate.Value
		input.DueDate = &t
	}

	t, err := h.taskService.CreateTask(ctx, input)
	if err != nil {
		return h.handleError(err).(api.CreateTaskRes), nil
	}

	return toAPITask(t), nil
}
