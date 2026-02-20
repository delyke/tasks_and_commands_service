package task

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestUpdateTask_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	newTitle := "Обновленный заголовок"

	input := UpdateTaskInput{
		TaskID:    existingTask.ID,
		Title:     &newTitle,
		UpdatedBy: userID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Update(s.ctx, mock.MatchedBy(func(t *task.Task) bool {
			return t.Title == newTitle
		})).
		Return(nil)

	s.historyRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(s.randomHistoryID(), nil)

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(newTitle, result.Title)
}

func (s *ServiceSuite) TestUpdateTask_MultipleFieldChanges() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	existingTask.Status = task.StatusTodo
	existingTask.Priority = task.PriorityLow

	newTitle := "Обновленный заголовок"
	newDescription := "Обновленное описание"
	newStatus := "in_progress"
	newPriority := "high"

	input := UpdateTaskInput{
		TaskID:      existingTask.ID,
		Title:       &newTitle,
		Description: &newDescription,
		Status:      &newStatus,
		Priority:    &newPriority,
		UpdatedBy:   userID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Update(s.ctx, mock.MatchedBy(func(t *task.Task) bool {
			return t.Title == newTitle &&
				t.Description == newDescription &&
				t.Status == task.StatusInProgress &&
				t.Priority == task.PriorityHigh
		})).
		Return(nil)

	s.historyRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(s.randomHistoryID(), nil).
		Times(4) // заголовок, описание, статус, приоритет

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(task.StatusInProgress, result.Status)
	s.Equal(task.PriorityHigh, result.Priority)
}

func (s *ServiceSuite) TestUpdateTask_WithAssigneeChange() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	newAssigneeID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)

	input := UpdateTaskInput{
		TaskID:     existingTask.ID,
		AssigneeID: &newAssigneeID,
		UpdatedBy:  userID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, newAssigneeID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Update(s.ctx, mock.MatchedBy(func(t *task.Task) bool {
			return t.AssigneeID != nil && *t.AssigneeID == newAssigneeID
		})).
		Return(nil)

	s.historyRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(s.randomHistoryID(), nil)

	s.cache.EXPECT().
		Delete(s.ctx, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.AssigneeID)
	s.Equal(newAssigneeID, *result.AssigneeID)
}

func (s *ServiceSuite) TestUpdateTask_TaskNotFound() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()
	newTitle := "Обновленный заголовок"

	input := UpdateTaskInput{
		TaskID:    taskID,
		Title:     &newTitle,
		UpdatedBy: userID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestUpdateTask_UserNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	creatorID := s.randomUserID()
	updaterID := s.randomUserID()
	existingTask := s.generateTask(teamID, creatorID)
	newTitle := "Обновленный заголовок"

	input := UpdateTaskInput{
		TaskID:    existingTask.ID,
		Title:     &newTitle,
		UpdatedBy: updaterID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, updaterID).
		Return(false, nil)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestUpdateTask_AssigneeNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	nonMemberAssignee := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)

	input := UpdateTaskInput{
		TaskID:     existingTask.ID,
		AssigneeID: &nonMemberAssignee,
		UpdatedBy:  userID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, nonMemberAssignee).
		Return(false, nil)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrInvalidInput)
}

func (s *ServiceSuite) TestUpdateTask_UpdateError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	newTitle := "Обновленный заголовок"
	expectedErr := errors.New("ошибка обновления")

	input := UpdateTaskInput{
		TaskID:    existingTask.ID,
		Title:     &newTitle,
		UpdatedBy: userID,
	}

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.taskRepo.EXPECT().
		Update(s.ctx, mock.Anything).
		Return(expectedErr)

	// Действие
	result, err := s.service.UpdateTask(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
