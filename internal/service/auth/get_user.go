package auth

import (
	"context"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// GetUser retrieves a user by ID.
func (s *Service) GetUser(ctx context.Context, id user.ID) (*user.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
