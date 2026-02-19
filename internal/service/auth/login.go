package auth

import (
	"context"
	"errors"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// LoginInput represents login input data.
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput represents login output data.
type LoginOutput struct {
	UserID user.ID
	Token  string
}

// Login authenticates a user.
func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	email, err := user.NewEmail(input.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate token
	token, err := s.tokenGen.Generate(u.ID, email.String(), s.config.AccessTokenTTL(), s.config.Issuer())
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		UserID: u.ID,
		Token:  token,
	}, nil
}
