package handler

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	"github.com/delyke/tasks_and_commands_service/internal/service/team"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// CreateTeam implements api.Handler.
func (h *Handler) CreateTeam(ctx context.Context, req *api.CreateTeamReq) (api.CreateTeamRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	description := ""
	if req.Description.IsSet() {
		description = req.Description.Value
	}

	t, err := h.teamService.CreateTeam(ctx, team.CreateTeamInput{
		Name:        req.Name,
		Description: description,
		CreatedBy:   user.ID(userID),
	})
	if err != nil {
		return h.handleError(err).(api.CreateTeamRes), nil
	}

	return toAPITeam(t), nil
}
