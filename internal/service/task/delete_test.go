package task

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
)

func (s *ServiceSuite) TestDeleteTask_Success() {
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

	s.taskRepo.EXPECT().
		Delete(s.ctx, existingTask.ID).
		Return(nil)

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	err := s.service.DeleteTask(s.ctx, existingTask.ID, userID)

	// Проверка
	s.NoError(err)
}

func (s *ServiceSuite) TestDeleteTask_TaskNotFound() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	err := s.service.DeleteTask(s.ctx, taskID, userID)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestDeleteTask_UserNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	deletedBy := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, deletedBy).
		Return(false, nil)

	// Действие
	err := s.service.DeleteTask(s.ctx, existingTask.ID, deletedBy)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestDeleteTask_MembershipCheckError() {
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
	err := s.service.DeleteTask(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestDeleteTask_DeleteError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	expectedErr := errors.New("ошибка удаления")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Delete(s.ctx, existingTask.ID).
		Return(expectedErr)

	// Действие
	err := s.service.DeleteTask(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Equal(expectedErr, err)
}
