package analytics

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestFindOrphanedAssignees_Success() {
	// Подготовка
	expectedOrphans := []*task.OrphanedTask{
		s.generateOrphanedTask(),
		s.generateOrphanedTask(),
	}

	s.analyticsRepo.EXPECT().
		FindOrphanedAssignees(s.ctx).
		Return(expectedOrphans, nil)

	// Действие
	result, err := s.service.FindOrphanedAssignees(s.ctx)

	// Проверка
	s.NoError(err)
	s.Len(result, 2)
	s.Equal(expectedOrphans, result)
}

func (s *ServiceSuite) TestFindOrphanedAssignees_NoOrphans() {
	// Подготовка
	s.analyticsRepo.EXPECT().
		FindOrphanedAssignees(s.ctx).
		Return([]*task.OrphanedTask{}, nil)

	// Действие
	result, err := s.service.FindOrphanedAssignees(s.ctx)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}

func (s *ServiceSuite) TestFindOrphanedAssignees_RepoError() {
	// Подготовка
	expectedErr := errors.New("ошибка базы данных")

	s.analyticsRepo.EXPECT().
		FindOrphanedAssignees(s.ctx).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.FindOrphanedAssignees(s.ctx)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
