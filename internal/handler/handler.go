package handler

import (
	"context"
	"errors"
	"github.com/delyke/tasks_and_commands_service/internal/service/auth"
	"github.com/delyke/tasks_and_commands_service/internal/service/task"
	"github.com/delyke/tasks_and_commands_service/internal/service/team"
	"net/http"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

// Handler implements the ogen Handler interface.
type Handler struct {
	authService *auth.Service
	teamService *team.Service
	taskService *task.Service
	userRepo    user.Repository
}

// NewHandler creates a new Handler.
func NewHandler(
	authService *auth.Service,
	teamService *team.Service,
	taskService *task.Service,
	userRepo user.Repository,
) *Handler {
	return &Handler{
		authService: authService,
		teamService: teamService,
		taskService: taskService,
		userRepo:    userRepo,
	}
}

// NewError implements api.Handler.
func (h *Handler) NewError(_ context.Context, err error) *api.GenericErrorStatusCode {
	return &api.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.GenericError{
			Code:    500,
			Message: err.Error(),
		},
	}
}

func (h *Handler) handleError(err error) any {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return &api.NotFoundError{Code: 404, Message: "not found"}
	case errors.Is(err, domain.ErrAlreadyExists):
		return &api.ConflictError{Code: 409, Message: "already exists"}
	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidCredentials):
		return &api.UnauthorizedError{Code: 401, Message: "invalid credentials"}
	case errors.Is(err, domain.ErrForbidden):
		return &api.ForbiddenError{Code: 403, Message: "forbidden"}
	case errors.Is(err, domain.ErrInvalidInput):
		return &api.BadRequestError{
			Code:    400,
			Message: err.Error(),
		}
	default:
		return &api.InternalServerError{Code: 500, Message: "internal server error"}
	}
}
