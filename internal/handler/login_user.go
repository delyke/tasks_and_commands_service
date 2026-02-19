package handler

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/service/auth"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// LoginUser implements api.Handler.
func (h *Handler) LoginUser(ctx context.Context, req *api.LoginUserReq) (api.LoginUserRes, error) {
	output, err := h.authService.Login(ctx, auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return h.handleError(err).(api.LoginUserRes), nil
	}

	u, err := h.authService.GetUser(ctx, output.UserID)
	if err != nil {
		return h.handleError(err).(api.LoginUserRes), nil
	}

	return &api.LoginUserOK{
		Token: output.Token,
		User:  toAPIUser(u),
	}, nil
}
