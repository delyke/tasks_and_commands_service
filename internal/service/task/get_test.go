package task

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
)

func (s *ServiceSuite) TestGetTask_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	expectedTask := s.generateTask(teamID, userID)

	s.taskRepo.EXPECT().
		GetByID(s.ctx, expectedTask.ID).
		Return(expectedTask, nil)

	// Действие
	result, err := s.service.GetTask(s.ctx, expectedTask.ID)

	// Проверка
	s.NoError(err)
	s.Equal(expectedTask, result)
}

func (s *ServiceSuite) TestGetTask_NotFound() {
	// Подготовка
	taskID := s.randomTaskID()

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.GetTask(s.ctx, taskID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestGetTask_RepoError() {
	// Подготовка
	taskID := s.randomTaskID()
	expectedErr := errors.New("ошибка базы данных")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.GetTask(s.ctx, taskID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
