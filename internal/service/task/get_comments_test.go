package task

import (
	"errors"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestGetTaskComments_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)

	expectedComments := []*task.Comment{
		s.generateComment(existingTask.ID, userID),
		s.generateComment(existingTask.ID, s.randomUserID()),
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.commentRepo.EXPECT().
		ListByTask(s.ctx, existingTask.ID).
		Return(expectedComments, nil)

	// Действие
	result, err := s.service.GetTaskComments(s.ctx, existingTask.ID, userID)

	// Проверка
	s.NoError(err)
	s.Len(result, 2)
}

func (s *ServiceSuite) TestGetTaskComments_TaskNotFound() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.GetTaskComments(s.ctx, taskID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestGetTaskComments_UserNotMember() {
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
	result, err := s.service.GetTaskComments(s.ctx, existingTask.ID, requesterID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestGetTaskComments_MembershipError() {
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
	result, err := s.service.GetTaskComments(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestGetTaskComments_ListError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	expectedErr := errors.New("ошибка получения комментариев")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.commentRepo.EXPECT().
		ListByTask(s.ctx, existingTask.ID).
		Return(nil, expectedErr)

	// Действие
	result, err := s.service.GetTaskComments(s.ctx, existingTask.ID, userID)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestGetTaskComments_EmptyComments() {
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

	s.commentRepo.EXPECT().
		ListByTask(s.ctx, existingTask.ID).
		Return([]*task.Comment{}, nil)

	// Действие
	result, err := s.service.GetTaskComments(s.ctx, existingTask.ID, userID)

	// Проверка
	s.NoError(err)
	s.Empty(result)
}
