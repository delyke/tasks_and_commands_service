package analytics

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestGetTopCreatorsByTeam_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	expectedCreators := []*task.TopCreator{
		s.generateTopCreator(teamID, 1),
		s.generateTopCreator(teamID, 2),
		s.generateTopCreator(teamID, 3),
	}

	s.analyticsRepo.EXPECT().
		GetTopCreatorsByTeam(s.ctx, 5).
		Return(expectedCreators, nil)

	// Действие
	result, err := s.service.GetTopCreatorsByTeam(s.ctx, 5)

	// Проверка
	s.NoError(err)
	s.Len(result, 3)
	s.Equal(expectedCreators, result)
}

func (s *ServiceSuite) TestGetTopCreatorsByTeam_DefaultLimit() {
	// Подготовка
	teamID := s.randomTeamID()
	expectedCreators := []*task.TopCreator{
		s.generateTopCreator(teamID, 1),
		s.generateTopCreator(teamID, 2),
		s.generateTopCreator(teamID, 3),
	}

	// limit <= 0 должен установиться в 3 по умолчанию
	s.analyticsRepo.EXPECT().
		GetTopCreatorsByTeam(s.ctx, 3).
		Return(expectedCreators, nil)

	// Действие
	result, err := s.service.GetTopCreatorsByTeam(s.ctx, 0)

	// Проверка
	s.NoError(err)
	s.Len(result, 3)
}

func (s *ServiceSuite) TestGetTopCreatorsByTeam_NegativeLimit() {
	// Подготовка
	teamID := s.randomTeamID()
	expectedCreators := []*task.TopCreator{
		s.generateTopCreator(teamID, 1),
	}

	// отрицательный limit должен установиться в 3 по умолчанию
	s.analyticsRepo.EXPECT().
		GetTopCreatorsByTeam(s.ctx, 3).
		Return(expectedCreators, nil)

	// Действие
	result, err := s.service.GetTopCreatorsByTeam(s.ctx, -5)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
}

func (s *ServiceSuite) TestGetTopCreatorsByTeam_EmptyResult() {
	// Подготовка
	s.analyticsRepo.EXPECT().
		GetTopCreatorsByTeam(s.ctx, 3).
		Return([]*task.TopCreator{}, nil)

	// Действие
	result, err := s.service.GetTopCreatorsByTeam(s.ctx, 3)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}

func (s *ServiceSuite) TestGetTopCreatorsByTeam_RepoError() {
	// Подготовка
	expectedErr := errors.New("ошибка базы данных")

	s.analyticsRepo.EXPECT().
		GetTopCreatorsByTeam(s.ctx, 5).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.GetTopCreatorsByTeam(s.ctx, 5)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestGetTopCreatorsByTeam_MultipleTeams() {
	// Подготовка
	teamID1 := s.randomTeamID()
	teamID2 := s.randomTeamID()
	expectedCreators := []*task.TopCreator{
		s.generateTopCreator(teamID1, 1),
		s.generateTopCreator(teamID1, 2),
		s.generateTopCreator(teamID2, 1),
		s.generateTopCreator(teamID2, 2),
	}

	s.analyticsRepo.EXPECT().
		GetTopCreatorsByTeam(s.ctx, 2).
		Return(expectedCreators, nil)

	// Действие
	result, err := s.service.GetTopCreatorsByTeam(s.ctx, 2)

	// Проверка
	s.NoError(err)
	s.Len(result, 4)
}
