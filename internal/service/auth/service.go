package auth

import (
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Config provides JWT configuration.
type Config interface {
	AccessTokenTTL() time.Duration
	Issuer() string
}

// TokenGenerator generates JWT tokens.
type TokenGenerator interface {
	Generate(userID user.ID, email string, ttl time.Duration, issuer string) (string, error)
}

// Service handles authentication operations.
type Service struct {
	userRepo user.Repository
	tokenGen TokenGenerator
	config   Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo user.Repository, tokenGen TokenGenerator, config Config) *Service {
	return &Service{
		userRepo: userRepo,
		tokenGen: tokenGen,
		config:   config,
	}
}
