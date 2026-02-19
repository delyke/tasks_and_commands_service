package middleware

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/service/jwt"
	"strings"

	api "github.com/delyke/tasks_and_commands_service/pkg/openapi/tasks/v1"
)

type contextKey string

const userIDContextKey contextKey = "user_id"

// SecurityHandler implements ogen security handler.
type SecurityHandler struct {
	jwtGenerator *jwt.Generator
}

// NewSecurityHandler creates a new SecurityHandler.
func NewSecurityHandler(jwtGenerator *jwt.Generator) *SecurityHandler {
	return &SecurityHandler{
		jwtGenerator: jwtGenerator,
	}
}

// HandleBearerAuth implements api.SecurityHandler.
func (h *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	token := t.Token
	if token == "" {
		return ctx, nil
	}

	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := h.jwtGenerator.Parse(token)
	if err != nil {
		return ctx, nil
	}

	return context.WithValue(ctx, userIDContextKey, uint64(claims.UserID)), nil
}

// GetUserIDFromContext extracts user ID from context.
func GetUserIDFromContext(ctx context.Context) uint64 {
	if v := ctx.Value(userIDContextKey); v != nil {
		if id, ok := v.(uint64); ok {
			return id
		}
	}
	return 0
}
