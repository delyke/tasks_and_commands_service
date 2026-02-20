package task

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestListTasks_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	teamIDVal := uint64(teamID)

	filter := &task.Filter{
		TeamID: &teamIDVal,
		Limit:  10,
		Offset: 0,
	}

	expectedTasks := []*task.Task{
		s.generateTask(teamID, userID),
		s.generateTask(teamID, userID),
	}

	s.cache.EXPECT().
		Get(s.ctx, mock.Anything, mock.Anything).
		Return(errors.New("кэш пуст"))

	s.taskRepo.EXPECT().
		List(s.ctx, filter).
		Return(expectedTasks, int64(2), nil)

	s.cache.EXPECT().
		Set(s.ctx, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Tasks, 2)
	s.Equal(int64(2), result.Total)
}

func (s *ServiceSuite) TestListTasks_FromCache() {
	// Подготовка
	teamID := s.randomTeamID()
	teamIDVal := uint64(teamID)

	filter := &task.Filter{
		TeamID: &teamIDVal,
		Limit:  10,
		Offset: 0,
	}

	// Кэш возвращает nil - данные найдены в кэше
	s.cache.EXPECT().
		Get(s.ctx, mock.Anything, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка - возвращает пустой результат (dest не заполнен моком)
	s.NoError(err)
	s.NotNil(result)
}

func (s *ServiceSuite) TestListTasks_WithStatusFilter() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	teamIDVal := uint64(teamID)
	status := task.StatusInProgress

	filter := &task.Filter{
		TeamID: &teamIDVal,
		Status: &status,
		Limit:  10,
		Offset: 0,
	}

	expectedTasks := []*task.Task{
		s.generateTask(teamID, userID),
	}
	expectedTasks[0].Status = task.StatusInProgress

	// Фильтрованные запросы не кэшируются
	s.taskRepo.EXPECT().
		List(s.ctx, filter).
		Return(expectedTasks, int64(1), nil)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Tasks, 1)
	s.Equal(task.StatusInProgress, result.Tasks[0].Status)
}

func (s *ServiceSuite) TestListTasks_WithAssigneeFilter() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	assigneeID := s.randomUserID()
	teamIDVal := uint64(teamID)

	filter := &task.Filter{
		TeamID:     &teamIDVal,
		AssigneeID: &assigneeID,
		Limit:      10,
		Offset:     0,
	}

	expectedTasks := []*task.Task{
		s.generateTask(teamID, userID),
	}
	expectedTasks[0].AssigneeID = &assigneeID

	// Фильтрованные запросы не кэшируются
	s.taskRepo.EXPECT().
		List(s.ctx, filter).
		Return(expectedTasks, int64(1), nil)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Tasks, 1)
}

func (s *ServiceSuite) TestListTasks_WithOffset() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	teamIDVal := uint64(teamID)

	filter := &task.Filter{
		TeamID: &teamIDVal,
		Limit:  10,
		Offset: 5, // ненулевой offset - не кэшируется
	}

	expectedTasks := []*task.Task{
		s.generateTask(teamID, userID),
	}

	// Запросы с offset не кэшируются
	s.taskRepo.EXPECT().
		List(s.ctx, filter).
		Return(expectedTasks, int64(6), nil)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Tasks, 1)
}

func (s *ServiceSuite) TestListTasks_EmptyResult() {
	// Подготовка
	teamID := s.randomTeamID()
	teamIDVal := uint64(teamID)

	filter := &task.Filter{
		TeamID: &teamIDVal,
		Limit:  10,
		Offset: 0,
	}

	s.cache.EXPECT().
		Get(s.ctx, mock.Anything, mock.Anything).
		Return(errors.New("кэш пуст"))

	s.taskRepo.EXPECT().
		List(s.ctx, filter).
		Return([]*task.Task{}, int64(0), nil)

	s.cache.EXPECT().
		Set(s.ctx, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Empty(result.Tasks)
	s.Equal(int64(0), result.Total)
}

func (s *ServiceSuite) TestListTasks_RepoError() {
	// Подготовка
	teamID := s.randomTeamID()
	teamIDVal := uint64(teamID)
	expectedErr := errors.New("ошибка базы данных")

	filter := &task.Filter{
		TeamID: &teamIDVal,
		Limit:  10,
		Offset: 0,
	}

	s.cache.EXPECT().
		Get(s.ctx, mock.Anything, mock.Anything).
		Return(errors.New("кэш пуст"))

	s.taskRepo.EXPECT().
		List(s.ctx, filter).
		Return(nil, int64(0), expectedErr)

	// Действие
	result, err := s.service.ListTasks(s.ctx, filter)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
