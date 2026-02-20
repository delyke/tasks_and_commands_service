package task

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	taskMocks "github.com/delyke/tasks_and_commands_service/internal/domain/task/mocks"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	teamMocks "github.com/delyke/tasks_and_commands_service/internal/domain/team/mocks"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	serviceMocks "github.com/delyke/tasks_and_commands_service/internal/service/task/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx         context.Context //nolint:containedctx
	taskRepo    *taskMocks.Repository
	historyRepo *taskMocks.HistoryRepository
	commentRepo *taskMocks.CommentRepository
	memberRepo  *teamMocks.MemberRepository
	cache       *serviceMocks.Cache
	service     *Service
	faker       *gofakeit.Faker
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.taskRepo = taskMocks.NewRepository(s.T())
	s.historyRepo = taskMocks.NewHistoryRepository(s.T())
	s.commentRepo = taskMocks.NewCommentRepository(s.T())
	s.memberRepo = teamMocks.NewMemberRepository(s.T())
	s.cache = serviceMocks.NewCache(s.T())
	s.service = NewTaskService(s.taskRepo, s.historyRepo, s.commentRepo, s.memberRepo, s.cache)
	s.faker = gofakeit.New(42)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

// Вспомогательные методы

func (s *ServiceSuite) randomUserID() user.ID {
	return user.ID(s.faker.Uint64())
}

func (s *ServiceSuite) randomTeamID() team.ID {
	return team.ID(s.faker.Uint64())
}

func (s *ServiceSuite) randomTaskID() task.ID {
	return task.ID(s.faker.Uint64())
}

func (s *ServiceSuite) randomCommentID() task.CommentID {
	return task.CommentID(s.faker.Uint64())
}

func (s *ServiceSuite) randomHistoryID() task.HistoryID {
	return task.HistoryID(s.faker.Uint64())
}

func (s *ServiceSuite) randomTitle() string {
	return s.faker.Sentence(5)
}

func (s *ServiceSuite) randomDescription() string {
	return s.faker.Paragraph(1, 3, 10, " ")
}

func (s *ServiceSuite) randomPriority() task.Priority {
	priorities := []task.Priority{task.PriorityLow, task.PriorityMedium, task.PriorityHigh}
	return priorities[s.faker.Number(0, len(priorities)-1)]
}

func (s *ServiceSuite) randomStatus() task.Status {
	statuses := []task.Status{task.StatusTodo, task.StatusInProgress, task.StatusDone}
	return statuses[s.faker.Number(0, len(statuses)-1)]
}

func (s *ServiceSuite) generateTask(teamID team.ID, createdBy user.ID) *task.Task {
	return &task.Task{
		ID:          s.randomTaskID(),
		TeamID:      teamID,
		Title:       s.randomTitle(),
		Description: s.randomDescription(),
		Status:      s.randomStatus(),
		Priority:    s.randomPriority(),
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (s *ServiceSuite) generateComment(taskID task.ID, userID user.ID) *task.Comment {
	return &task.Comment{
		ID:        s.randomCommentID(),
		TaskID:    taskID,
		UserID:    userID,
		Content:   s.faker.Sentence(10),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *ServiceSuite) generateHistory(taskID task.ID, changedBy user.ID) *task.History {
	oldVal := s.faker.Word()
	newVal := s.faker.Word()
	return &task.History{
		ID:        s.randomHistoryID(),
		TaskID:    taskID,
		ChangedBy: changedBy,
		FieldName: "status",
		OldValue:  &oldVal,
		NewValue:  &newVal,
		ChangedAt: time.Now(),
	}
}
