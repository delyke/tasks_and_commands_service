package task

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
)

func (s *ServiceSuite) TestVerifyTaskAccess_Success() {
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

	// Действие
	result, err := s.service.VerifyTaskAccess(s.ctx, existingTask.ID, userID)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(existingTask.ID, result.ID)
}

func (s *ServiceSuite) TestVerifyTaskAccess_TaskNotFound() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.VerifyTaskAccess(s.ctx, taskID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestVerifyTaskAccess_UserNotMember() {
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
	result, err := s.service.VerifyTaskAccess(s.ctx, existingTask.ID, requesterID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestVerifyTaskAccess_MembershipError() {
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
	result, err := s.service.VerifyTaskAccess(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestVerifyTaskAccess_GenericRepoError() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()
	expectedErr := errors.New("ошибка базы данных")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.VerifyTaskAccess(s.ctx, taskID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
