package analytics

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestGetTeamStats_Success() {
	// Подготовка
	expectedStats := []*task.TeamStats{
		s.generateTeamStats(),
		s.generateTeamStats(),
		s.generateTeamStats(),
	}

	s.analyticsRepo.EXPECT().
		GetTeamStats(s.ctx).
		Return(expectedStats, nil)

	// Действие
	result, err := s.service.GetTeamStats(s.ctx)

	// Проверка
	s.NoError(err)
	s.Len(result, 3)
	s.Equal(expectedStats, result)
}

func (s *ServiceSuite) TestGetTeamStats_EmptyResult() {
	// Подготовка
	s.analyticsRepo.EXPECT().
		GetTeamStats(s.ctx).
		Return([]*task.TeamStats{}, nil)

	// Действие
	result, err := s.service.GetTeamStats(s.ctx)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}

func (s *ServiceSuite) TestGetTeamStats_RepoError() {
	// Подготовка
	expectedErr := errors.New("ошибка базы данных")

	s.analyticsRepo.EXPECT().
		GetTeamStats(s.ctx).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.GetTeamStats(s.ctx)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
