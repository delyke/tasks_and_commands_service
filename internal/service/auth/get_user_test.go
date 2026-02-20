package auth

import (
	"github.com/delyke/tasks_and_commands_service/internal/domain"
)

func (s *ServiceSuite) TestGetUser_Success() {
	// Подготовка
	expectedUser := s.generateUser(s.randomPassword())

	s.userRepo.EXPECT().
		GetByID(s.ctx, expectedUser.ID).
		Return(expectedUser, nil)

	// Действие
	result, err := s.service.GetUser(s.ctx, expectedUser.ID)

	// Проверка
	s.NoError(err)
	s.Equal(expectedUser, result)
}

func (s *ServiceSuite) TestGetUser_NotFound() {
	// Подготовка
	userID := s.randomUserID()

	s.userRepo.EXPECT().
		GetByID(s.ctx, userID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.GetUser(s.ctx, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}
