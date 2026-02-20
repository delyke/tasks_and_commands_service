package auth

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	userMocks "github.com/delyke/tasks_and_commands_service/internal/domain/user/mocks"
	authMocks "github.com/delyke/tasks_and_commands_service/internal/service/auth/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx      context.Context //nolint:containedctx
	userRepo *userMocks.Repository
	tokenGen *authMocks.TokenGenerator
	config   *authMocks.Config
	service  *Service
	faker    *gofakeit.Faker
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = userMocks.NewRepository(s.T())
	s.tokenGen = authMocks.NewTokenGenerator(s.T())
	s.config = authMocks.NewConfig(s.T())
	s.service = NewAuthService(s.userRepo, s.tokenGen, s.config)
	s.faker = gofakeit.New(42)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

// Вспомогательные методы

func (s *ServiceSuite) randomUserID() user.ID {
	return user.ID(s.faker.Uint64())
}

func (s *ServiceSuite) randomEmail() string {
	return s.faker.Email()
}

func (s *ServiceSuite) randomPassword() string {
	return s.faker.Password(true, true, true, true, false, 16)
}

func (s *ServiceSuite) randomName() string {
	return s.faker.Name()
}

func (s *ServiceSuite) randomToken() string {
	return s.faker.UUID()
}

func (s *ServiceSuite) hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (s *ServiceSuite) generateUser(password string) *user.User {
	email, _ := user.NewEmail(s.randomEmail())
	return &user.User{
		ID:           s.randomUserID(),
		Email:        email,
		PasswordHash: s.hashPassword(password),
		Name:         s.randomName(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (s *ServiceSuite) setupConfigExpectations() {
	s.config.EXPECT().AccessTokenTTL().Return(time.Hour).Maybe()
	s.config.EXPECT().Issuer().Return("test-issuer").Maybe()
}
