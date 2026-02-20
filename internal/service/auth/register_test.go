package auth

import (
	"errors"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

func (s *ServiceSuite) TestRegister_Success() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)
	password := s.randomPassword()
	name := s.randomName()
	userID := s.randomUserID()
	expectedToken := s.randomToken()

	input := RegisterInput{
		Email:    emailStr,
		Password: password,
		Name:     name,
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(nil, domain.ErrNotFound)

	s.userRepo.EXPECT().
		Create(s.ctx, mock.MatchedBy(func(u *user.User) bool {
			return u.Email == email && u.Name == name && u.PasswordHash != ""
		})).
		Return(userID, nil)

	s.config.EXPECT().AccessTokenTTL().Return(time.Hour)
	s.config.EXPECT().Issuer().Return("test-issuer")

	s.tokenGen.EXPECT().
		Generate(userID, emailStr, time.Hour, "test-issuer").
		Return(expectedToken, nil)

	// Действие
	result, err := s.service.Register(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(userID, result.UserID)
	s.Equal(expectedToken, result.Token)
}

func (s *ServiceSuite) TestRegister_InvalidEmail() {
	// Подготовка
	input := RegisterInput{
		Email:    "invalid-email",
		Password: s.randomPassword(),
		Name:     s.randomName(),
	}

	// Действие
	result, err := s.service.Register(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, user.ErrInvalidEmail)
}

func (s *ServiceSuite) TestRegister_UserAlreadyExists() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)
	existingUser := s.generateUser(s.randomPassword())
	existingUser.Email = email

	input := RegisterInput{
		Email:    emailStr,
		Password: s.randomPassword(),
		Name:     s.randomName(),
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(existingUser, nil)

	// Действие
	result, err := s.service.Register(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrAlreadyExists)
}

func (s *ServiceSuite) TestRegister_GetByEmailError() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)
	expectedErr := errors.New("ошибка базы данных")

	input := RegisterInput{
		Email:    emailStr,
		Password: s.randomPassword(),
		Name:     s.randomName(),
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.Register(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestRegister_CreateUserError() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)
	expectedErr := errors.New("ошибка создания")

	input := RegisterInput{
		Email:    emailStr,
		Password: s.randomPassword(),
		Name:     s.randomName(),
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(nil, domain.ErrNotFound)

	s.userRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(user.ID(0), expectedErr)

	// Действие
	result, err := s.service.Register(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestRegister_TokenGenerationError() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)
	userID := s.randomUserID()
	expectedErr := errors.New("ошибка генерации токена")

	input := RegisterInput{
		Email:    emailStr,
		Password: s.randomPassword(),
		Name:     s.randomName(),
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(nil, domain.ErrNotFound)

	s.userRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(userID, nil)

	s.config.EXPECT().AccessTokenTTL().Return(time.Hour)
	s.config.EXPECT().Issuer().Return("test-issuer")

	s.tokenGen.EXPECT().
		Generate(userID, emailStr, time.Hour, "test-issuer").
		Return("", expectedErr)

	// Действие
	result, err := s.service.Register(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
