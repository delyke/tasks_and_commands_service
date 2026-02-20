package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// RegisterInput represents registration input data.
type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

// RegisterOutput represents registration output data.
type RegisterOutput struct {
	UserID user.ID
	Token  string
}

// Register creates a new user account.
func (s *Service) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	email, err := user.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	_, err = s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	u := user.New(email, string(hashedPassword), input.Name)
	userID, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.tokenGen.Generate(userID, email.String(), s.config.AccessTokenTTL(), s.config.Issuer())
	if err != nil {
		return nil, err
	}

	return &RegisterOutput{
		UserID: userID,
		Token:  token,
	}, nil
}
