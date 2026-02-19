package handler

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// ListTasks implements api.Handler.
func (h *Handler) ListTasks(ctx context.Context, params api.ListTasksParams) (api.ListTasksRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	filter := task.NewFilter()

	if params.TeamID.IsSet() {
		teamID := uint64(params.TeamID.Value) //nolint:gosec // safe: API ID is validated
		filter = filter.WithTeamID(teamID)
	}

	if params.Status.IsSet() {
		status, err := task.NewStatus(string(params.Status.Value))
		if err == nil {
			filter = filter.WithStatus(status)
		}
	}

	if params.AssigneeID.IsSet() {
		filter = filter.WithAssignee(user.ID(params.AssigneeID.Value)) //nolint:gosec // safe: API ID is validated
	}

	page := params.Page.Or(1)
	perPage := params.PerPage.Or(20)
	offset := (page - 1) * perPage
	filter = filter.WithPagination(int(perPage), int(offset))

	result, err := h.taskService.ListTasks(ctx, filter)
	if err != nil {
		return h.handleError(err).(api.ListTasksRes), nil
	}

	tasks := make([]api.Task, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		tasks = append(tasks, *toAPITask(t))
	}

	totalPages := int32((result.Total + int64(perPage) - 1) / int64(perPage)) //nolint:gosec // safe: pagination values are bounded

	return &api.ListTasksOK{
		Tasks: tasks,
		Pagination: api.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      result.Total,
			TotalPages: totalPages,
		},
	}, nil
}
