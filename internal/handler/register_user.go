package handler

import (
	"context"

	"github.com/delyke/tasks_and_commands_service/internal/service/auth"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// RegisterUser implements api.Handler.
func (h *Handler) RegisterUser(ctx context.Context, req *api.RegisterUserReq) (api.RegisterUserRes, error) {
	output, err := h.authService.Register(ctx, auth.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Username,
	})
	if err != nil {
		return h.handleError(err).(api.RegisterUserRes), nil
	}

	u, err := h.authService.GetUser(ctx, output.UserID)
	if err != nil {
		return h.handleError(err).(api.RegisterUserRes), nil
	}

	return &api.RegisterUserOK{
		Token: output.Token,
		User:  toAPIUser(u),
	}, nil
}
