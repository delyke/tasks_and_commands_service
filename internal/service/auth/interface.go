package auth

import (
	"context"
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error)
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	GetUser(ctx context.Context, id user.ID) (*user.User, error)
}

// Config provides JWT configuration.
type Config interface {
	AccessTokenTTL() time.Duration
	Issuer() string
}

// TokenGenerator generates JWT tokens.
type TokenGenerator interface {
	Generate(userID user.ID, email string, ttl time.Duration, issuer string) (string, error)
}
