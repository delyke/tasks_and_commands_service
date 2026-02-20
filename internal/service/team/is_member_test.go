package team

import (
	"errors"
)

func (s *ServiceSuite) TestIsMember_True() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	// Действие
	result, err := s.service.IsMember(s.ctx, teamID, userID)

	// Проверка
	s.NoError(err)
	s.True(result)
}

func (s *ServiceSuite) TestIsMember_False() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(false, nil)

	// Действие
	result, err := s.service.IsMember(s.ctx, teamID, userID)

	// Проверка
	s.NoError(err)
	s.False(result)
}

func (s *ServiceSuite) TestIsMember_Error() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	expectedErr := errors.New("ошибка базы данных")

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(false, expectedErr)

	// Действие
	result, err := s.service.IsMember(s.ctx, teamID, userID)

	// Проверка
	s.Error(err)
	s.False(result)
	s.Equal(expectedErr, err)
}
