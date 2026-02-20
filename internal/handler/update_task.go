package handler

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	task2 "github.com/delyke/tasks_and_commands_service/internal/service/task"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// UpdateTask implements api.Handler.
func (h *Handler) UpdateTask(ctx context.Context, req *api.UpdateTaskReq, params api.UpdateTaskParams) (api.UpdateTaskRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	input := task2.UpdateTaskInput{
		TaskID:    task.ID(params.ID), //nolint:gosec // safe: API ID is validated
		UpdatedBy: user.ID(userID),    //nolint:gosec // safe: API ID is validated
	}

	if req.Title.IsSet() {
		input.Title = &req.Title.Value
	}
	if req.Description.IsSet() {
		input.Description = &req.Description.Value
	}
	if req.Status.IsSet() {
		s := string(req.Status.Value)
		input.Status = &s
	}
	if req.Priority.IsSet() {
		p := string(req.Priority.Value)
		input.Priority = &p
	}
	if req.AssigneeID.IsSet() && !req.AssigneeID.Null {
		id := user.ID(req.AssigneeID.Value) //nolint:gosec // safe: API ID is validated
		input.AssigneeID = &id
	}
	if req.DueDate.IsSet() && !req.DueDate.Null {
		input.DueDate = &req.DueDate.Value
	}

	t, err := h.taskService.UpdateTask(ctx, input)
	if err != nil {
		return h.handleError(err).(api.UpdateTaskRes), nil
	}

	return toAPITask(t), nil
}
