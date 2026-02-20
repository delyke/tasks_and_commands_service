package auth

import (
	"errors"
	"time"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

func (s *ServiceSuite) TestLogin_Success() {
	// Подготовка
	password := s.randomPassword()
	existingUser := s.generateUser(password)
	emailStr := existingUser.Email.String()
	expectedToken := s.randomToken()

	input := LoginInput{
		Email:    emailStr,
		Password: password,
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, existingUser.Email).
		Return(existingUser, nil)

	s.config.EXPECT().AccessTokenTTL().Return(time.Hour)
	s.config.EXPECT().Issuer().Return("test-issuer")

	s.tokenGen.EXPECT().
		Generate(existingUser.ID, emailStr, time.Hour, "test-issuer").
		Return(expectedToken, nil)

	// Действие
	result, err := s.service.Login(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(existingUser.ID, result.UserID)
	s.Equal(expectedToken, result.Token)
}

func (s *ServiceSuite) TestLogin_InvalidEmail() {
	// Подготовка
	input := LoginInput{
		Email:    "invalid-email",
		Password: s.randomPassword(),
	}

	// Действие
	result, err := s.service.Login(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrInvalidCredentials)
}

func (s *ServiceSuite) TestLogin_UserNotFound() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)

	input := LoginInput{
		Email:    emailStr,
		Password: s.randomPassword(),
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.Login(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrInvalidCredentials)
}

func (s *ServiceSuite) TestLogin_GetByEmailError() {
	// Подготовка
	emailStr := s.randomEmail()
	email, _ := user.NewEmail(emailStr)
	expectedErr := errors.New("ошибка базы данных")

	input := LoginInput{
		Email:    emailStr,
		Password: s.randomPassword(),
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, email).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.Login(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestLogin_WrongPassword() {
	// Подготовка
	password := s.randomPassword()
	existingUser := s.generateUser(password)
	emailStr := existingUser.Email.String()

	input := LoginInput{
		Email:    emailStr,
		Password: "wrong-password",
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, existingUser.Email).
		Return(existingUser, nil)

	// Действие
	result, err := s.service.Login(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrInvalidCredentials)
}

func (s *ServiceSuite) TestLogin_TokenGenerationError() {
	// Подготовка
	password := s.randomPassword()
	existingUser := s.generateUser(password)
	emailStr := existingUser.Email.String()
	expectedErr := errors.New("ошибка генерации токена")

	input := LoginInput{
		Email:    emailStr,
		Password: password,
	}

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, existingUser.Email).
		Return(existingUser, nil)

	s.config.EXPECT().AccessTokenTTL().Return(time.Hour)
	s.config.EXPECT().Issuer().Return("test-issuer")

	s.tokenGen.EXPECT().
		Generate(existingUser.ID, emailStr, time.Hour, "test-issuer").
		Return("", expectedErr)

	// Действие
	result, err := s.service.Login(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
