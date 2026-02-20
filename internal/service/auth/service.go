package auth

import (
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Service handles authentication operations.
type Service struct {
	userRepo user.Repository
	tokenGen TokenGenerator
	config   Config
}

// Ensure Service implements AuthService.
var _ AuthService = (*Service)(nil)

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo user.Repository, tokenGen TokenGenerator, config Config) *Service {
	return &Service{
		userRepo: userRepo,
		tokenGen: tokenGen,
		config:   config,
	}
}
