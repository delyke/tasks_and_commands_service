package handler

import (
	"context"
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"github.com/delyke/tasks_and_commands_service/internal/handler/middleware"
	team2 "github.com/delyke/tasks_and_commands_service/internal/service/team"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// InviteUserToTeam implements api.Handler.
func (h *Handler) InviteUserToTeam(ctx context.Context, req *api.InviteUserToTeamReq, params api.InviteUserToTeamParams) (api.InviteUserToTeamRes, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	if userID == 0 {
		return &api.UnauthorizedError{Code: 401, Message: "unauthorized"}, nil
	}

	// Get invitee user
	invitee, err := h.userRepo.GetByID(ctx, user.ID(req.UserID)) //nolint:gosec // safe: API ID is validated
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &api.NotFoundError{Code: 404, Message: "user not found"}, nil
		}
		return h.handleError(err).(api.InviteUserToTeamRes), nil
	}

	role := "member"
	if req.Role.IsSet() {
		role = string(req.Role.Value)
	}

	err = h.teamService.InviteMember(ctx, team2.InviteInput{
		TeamID:       team.ID(params.ID), //nolint:gosec // safe: API ID is validated
		InviterID:    user.ID(userID),    //nolint:gosec // safe: API ID is validated
		InviteeEmail: invitee.Email.String(),
		Role:         role,
	})
	if err != nil {
		return h.handleError(err).(api.InviteUserToTeamRes), nil
	}

	member, err := h.teamService.GetMember(ctx, team.ID(params.ID), user.ID(req.UserID)) //nolint:gosec // safe: API IDs are validated
	if err != nil {
		return h.handleError(err).(api.InviteUserToTeamRes), nil
	}

	return &api.TeamMember{
		UserID:   int64(member.UserID), //nolint:gosec // safe: domain ID is always positive
		TeamID:   int64(member.TeamID), //nolint:gosec // safe: domain ID is always positive
		Role:     api.TeamMemberRole(member.Role),
		JoinedAt: member.JoinedAt,
		User:     toAPIUser(invitee),
	}, nil
}
