package team

import (
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

func (s *ServiceSuite) TestGetTeam_Success() {
	// Подготовка
	userID := s.randomUserID()
	expectedTeam := s.generateTeam(userID)

	s.teamRepo.EXPECT().
		GetByID(s.ctx, expectedTeam.ID).
		Return(expectedTeam, nil)

	// Действие
	result, err := s.service.GetTeam(s.ctx, expectedTeam.ID)

	// Проверка
	s.NoError(err)
	s.Equal(expectedTeam, result)
}

func (s *ServiceSuite) TestGetTeam_NotFound() {
	// Подготовка
	teamID := s.randomTeamID()

	s.teamRepo.EXPECT().
		GetByID(s.ctx, teamID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.GetTeam(s.ctx, teamID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestGetMember_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	expectedMember := s.generateMember(teamID, userID, team.RoleMember)

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, userID).
		Return(expectedMember, nil)

	// Действие
	result, err := s.service.GetMember(s.ctx, teamID, userID)

	// Проверка
	s.NoError(err)
	s.Equal(expectedMember, result)
}

func (s *ServiceSuite) TestGetMember_NotFound() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, userID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.GetMember(s.ctx, teamID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}
