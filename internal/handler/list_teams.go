package handler

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// ListTeams implements api.Handler.
func (h *Handler) ListTeams(ctx context.Context, params api.ListTeamsParams) (api.ListTeamsRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	teams, err := h.teamService.ListUserTeams(ctx, user.ID(userID))
	if err != nil {
		return h.handleError(err).(api.ListTeamsRes), nil
	}

	apiTeams := make([]api.Team, 0, len(teams))
	for _, t := range teams {
		apiTeams = append(apiTeams, *toAPITeam(t))
	}

	page := params.Page.Or(1)
	perPage := params.PerPage.Or(20)
	total := int64(len(teams))
	totalPages := int32((total + int64(perPage) - 1) / int64(perPage)) //nolint:gosec // safe: pagination values are bounded

	return &api.ListTeamsOK{
		Teams: apiTeams,
		Pagination: api.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
