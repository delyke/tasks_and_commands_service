package task

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
)

func (s *ServiceSuite) TestAddComment_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	commentContent := s.faker.Sentence(10)
	commentID := s.randomCommentID()

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.commentRepo.EXPECT().
		Create(s.ctx, mock.MatchedBy(func(c *task.Comment) bool {
			return c.TaskID == existingTask.ID &&
				c.UserID == userID &&
				c.Content == commentContent
		})).
		Return(commentID, nil)

	// Действие
	result, err := s.service.AddComment(s.ctx, existingTask.ID, userID, commentContent)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(commentID, result.ID)
	s.Equal(existingTask.ID, result.TaskID)
	s.Equal(userID, result.UserID)
	s.Equal(commentContent, result.Content)
}

func (s *ServiceSuite) TestAddComment_TaskNotFound() {
	// Подготовка
	taskID := s.randomTaskID()
	userID := s.randomUserID()
	commentContent := s.faker.Sentence(10)

	s.taskRepo.EXPECT().
		GetByID(s.ctx, taskID).
		Return(nil, domain.ErrNotFound)

	// Действие
	result, err := s.service.AddComment(s.ctx, taskID, userID, commentContent)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestAddComment_UserNotMember() {
	// Подготовка
	teamID := s.randomTeamID()
	creatorID := s.randomUserID()
	commenterID := s.randomUserID()
	existingTask := s.generateTask(teamID, creatorID)
	commentContent := s.faker.Sentence(10)

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, commenterID).
		Return(false, nil)

	// Действие
	result, err := s.service.AddComment(s.ctx, existingTask.ID, commenterID, commentContent)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestAddComment_MembershipError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	commentContent := s.faker.Sentence(10)
	expectedErr := errors.New("ошибка базы данных")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(false, expectedErr)

	// Действие
	result, err := s.service.AddComment(s.ctx, existingTask.ID, userID, commentContent)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestAddComment_CreateError() {
	// Подготовка
	teamID := s.randomTeamID()
	userID := s.randomUserID()
	existingTask := s.generateTask(teamID, userID)
	commentContent := s.faker.Sentence(10)
	expectedErr := errors.New("ошибка создания комментария")

	s.taskRepo.EXPECT().
		GetByID(s.ctx, existingTask.ID).
		Return(existingTask, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, userID).
		Return(true, nil)

	s.commentRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(task.CommentID(0), expectedErr)

	// Действие
	result, err := s.service.AddComment(s.ctx, existingTask.ID, userID, commentContent)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
