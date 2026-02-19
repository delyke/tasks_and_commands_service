package handler

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// GetTaskHistory implements api.Handler.
func (h *Handler) GetTaskHistory(ctx context.Context, params api.GetTaskHistoryParams) (api.GetTaskHistoryRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	history, err := h.taskService.GetTaskHistory(ctx, task.ID(params.ID), user.ID(userID)) //nolint:gosec // safe: API IDs are validated
	if err != nil {
		return h.handleError(err).(api.GetTaskHistoryRes), nil
	}

	entries := make([]api.TaskHistoryEntry, 0, len(history))
	for _, histEntry := range history {
		changer, _ := h.userRepo.GetByID(ctx, histEntry.ChangedBy) //nolint:errcheck,gosec // optional: changer may not exist
		entry := api.TaskHistoryEntry{
			ID:        int64(histEntry.ID),        //nolint:gosec // safe: domain ID is always positive
			TaskID:    int64(histEntry.TaskID),    //nolint:gosec // safe: domain ID is always positive
			ChangedBy: int64(histEntry.ChangedBy), //nolint:gosec // safe: domain ID is always positive
			FieldName: histEntry.FieldName,
			ChangedAt: histEntry.ChangedAt,
		}
		if histEntry.OldValue != nil {
			entry.OldValue.SetTo(*histEntry.OldValue)
		}
		if histEntry.NewValue != nil {
			entry.NewValue.SetTo(*histEntry.NewValue)
		}
		if changer != nil {
			entry.Changer = toAPIUser(changer)
		}
		entries = append(entries, entry)
	}

	return &api.GetTaskHistoryOK{
		History: entries,
		Pagination: api.Pagination{
			Page:       1,
			PerPage:    int32(len(entries)), //nolint:gosec // safe: history entries are bounded
			Total:      int64(len(entries)),
			TotalPages: 1,
		},
	}, nil
}
