package team

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

func (s *ServiceSuite) TestListUserTeams_Success() {
	// Подготовка
	userID := s.randomUserID()
	expectedTeams := []*team.Team{
		s.generateTeam(userID),
		s.generateTeam(userID),
		s.generateTeam(s.randomUserID()),
	}

	s.teamRepo.EXPECT().
		ListByUser(s.ctx, userID).
		Return(expectedTeams, nil)

	// Действие
	result, err := s.service.ListUserTeams(s.ctx, userID)

	// Проверка
	s.NoError(err)
	s.Equal(expectedTeams, result)
	s.Len(result, 3)
}

func (s *ServiceSuite) TestListUserTeams_Empty() {
	// Подготовка
	userID := s.randomUserID()

	s.teamRepo.EXPECT().
		ListByUser(s.ctx, userID).
		Return([]*team.Team{}, nil)

	// Действие
	result, err := s.service.ListUserTeams(s.ctx, userID)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}

func (s *ServiceSuite) TestListUserTeams_Error() {
	// Подготовка
	userID := s.randomUserID()
	expectedErr := errors.New("ошибка базы данных")

	s.teamRepo.EXPECT().
		ListByUser(s.ctx, userID).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.ListUserTeams(s.ctx, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestListTeamMembers_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	expectedMembers := []*team.Member{
		s.generateMember(teamID, s.randomUserID(), team.RoleOwner),
		s.generateMember(teamID, s.randomUserID(), team.RoleAdmin),
		s.generateMember(teamID, s.randomUserID(), team.RoleMember),
	}

	s.memberRepo.EXPECT().
		ListByTeam(s.ctx, teamID).
		Return(expectedMembers, nil)

	// Действие
	result, err := s.service.ListTeamMembers(s.ctx, teamID)

	// Проверка
	s.NoError(err)
	s.Equal(expectedMembers, result)
	s.Len(result, 3)
}

func (s *ServiceSuite) TestListTeamMembers_Empty() {
	// Подготовка
	teamID := s.randomTeamID()

	s.memberRepo.EXPECT().
		ListByTeam(s.ctx, teamID).
		Return([]*team.Member{}, nil)

	// Действие
	result, err := s.service.ListTeamMembers(s.ctx, teamID)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}

func (s *ServiceSuite) TestListTeamMembers_Error() {
	// Подготовка
	teamID := s.randomTeamID()
	expectedErr := errors.New("ошибка базы данных")

	s.memberRepo.EXPECT().
		ListByTeam(s.ctx, teamID).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.ListTeamMembers(s.ctx, teamID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
