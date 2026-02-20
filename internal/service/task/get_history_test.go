package task

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestGetTaskHistory_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)

	expectedHistory := []*task.History{
		s.generateHistory(existingTask.ID, userID),
		s.generateHistory(existingTask.ID, userID),
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.historyRepo.EXPECT().
		ListByTask(s.ctx, existingTask.ID).
		Return(expectedHistory, nil)

	// Действие
	result, err := s.service.GetTaskHistory(s.ctx, existingTask.ID, userID)

	// Проверка
	s.NoError(err)
	s.Len(result, 2)
}

func (s *ServiceSuite) TestGetTaskHistory_TaskNotFound() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.GetTaskHistory(s.ctx, taskID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestGetTaskHistory_UserNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	creatorID := s.randomUserID()
	requesterID := s.randomUserID()
	existingTask := s.generateTask(teamID, creatorID)

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, requesterID).
		Return(false, nil)

	// Действие
	result, err := s.service.GetTaskHistory(s.ctx, existingTask.ID, requesterID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestGetTaskHistory_MembershipError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	expectedErr := errors.New("ошибка базы данных")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(false, expectedErr)

	// Действие
	result, err := s.service.GetTaskHistory(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestGetTaskHistory_ListError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	expectedErr := errors.New("ошибка получения истории")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.historyRepo.EXPECT().
		ListByTask(s.ctx, existingTask.ID).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.GetTaskHistory(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestGetTaskHistory_EmptyHistory() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.historyRepo.EXPECT().
		ListByTask(s.ctx, existingTask.ID).
		Return([]*task.History{}, nil)

	// Действие
	result, err := s.service.GetTaskHistory(s.ctx, existingTask.ID, userID)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}
