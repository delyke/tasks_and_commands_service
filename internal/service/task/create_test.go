package task

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestCreateTask_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	taskID := s.randomTaskID()

	input := CreateTaskInput{
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Priority:    "high",
		CreatedBy:   userID,
	}

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Create(s.ctx, mock.MatchedBy(func(t *task.Task) bool {
			return t.TeamID == teamID &&
				t.Title == input.Title &&
				t.Priority == task.PriorityHigh &&
				t.CreatedBy == userID
		})).
		Return(taskID, nil)

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.CreateTask(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(taskID, result.ID)
	s.Equal(input.Title, result.Title)
	s.Equal(task.PriorityHigh, result.Priority)
}

func (s *ServiceSuite) TestCreateTask_WithAssignee() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	assigneeID := s.randomUserID()
	taskID := s.randomTaskID()

	input := CreateTaskInput{
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Priority:    "medium",
		AssigneeID:  &assigneeID,
		CreatedBy:   userID,
	}

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, assigneeID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Create(s.ctx, mock.MatchedBy(func(t *task.Task) bool {
			return t.AssigneeID != nil && *t.AssigneeID == assigneeID
		})).
		Return(taskID, nil)

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.CreateTask(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.AssigneeID)
	s.Equal(assigneeID, *result.AssigneeID)
}

func (s *ServiceSuite) TestCreateTask_UserNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()

	input := CreateTaskInput{
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Priority:    "medium",
		CreatedBy:   userID,
	}

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(false, nil)

	// Действие
	result, err := s.service.CreateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestCreateTask_AssigneeNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	assigneeID := s.randomUserID()

	input := CreateTaskInput{
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Priority:    "medium",
		AssigneeID:  &assigneeID,
		CreatedBy:   userID,
	}

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, assigneeID).
		Return(false, nil)

	// Действие
	result, err := s.service.CreateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrInvalidInput)
}

func (s *ServiceSuite) TestCreateTask_InvalidPriorityDefaultsToMedium() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	taskID := s.randomTaskID()

	input := CreateTaskInput{
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Priority:    "invalid-priority",
		CreatedBy:   userID,
	}

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Create(s.ctx, mock.MatchedBy(func(t *task.Task) bool {
			return t.Priority == task.PriorityMedium // по умолчанию medium
		})).
		Return(taskID, nil)

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.CreateTask(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(task.PriorityMedium, result.Priority)
}

func (s *ServiceSuite) TestCreateTask_RepoError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	expectedErr := errors.New("ошибка базы данных")

	input := CreateTaskInput{
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Priority:    "medium",
		CreatedBy:   userID,
	}

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(task.ID(0), expectedErr)

	// Действие
	result, err := s.service.CreateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
