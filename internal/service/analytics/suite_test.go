package analytics

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/delyke/tasks_and_commands_service/internal/domain/task"
	taskMocks "github.com/delyke/tasks_and_commands_service/internal/domain/task/mocks"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

type ServiceSuite struct {
	suite.Suite
	ctx           context.Context //nolint:containedctx
	analyticsRepo *taskMocks.AnalyticsRepository
	service       *Service
	faker         *gofakeit.Faker
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.analyticsRepo = taskMocks.NewAnalyticsRepository(s.T())
	s.service = NewAnalyticsService(s.analyticsRepo)
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

func (s *ServiceSuite) randomTeamName() string {
	return s.faker.Company()
}

func (s *ServiceSuite) randomUserName() string {
	return s.faker.Name()
}

func (s *ServiceSuite) randomTaskTitle() string {
	return s.faker.Sentence(5)
}

func (s *ServiceSuite) generateTeamStats() *task.TeamStats {
	return &task.TeamStats{
		TeamID:          s.randomTeamID(),
		TeamName:        s.randomTeamName(),
		MemberCount:     int64(s.faker.Number(1, 50)),
		DoneTasksLast7d: int64(s.faker.Number(0, 100)),
	}
}

func (s *ServiceSuite) generateTopCreator(teamID team.ID, rank int) *task.TopCreator {
	return &task.TopCreator{
		TeamID:    teamID,
		TeamName:  s.randomTeamName(),
		UserID:    s.randomUserID(),
		UserName:  s.randomUserName(),
		TaskCount: int64(s.faker.Number(1, 50)),
		Rank:      rank,
	}
}

func (s *ServiceSuite) generateOrphanedTask() *task.OrphanedTask {
	return &task.OrphanedTask{
		TaskID:       s.randomTaskID(),
		TaskTitle:    s.randomTaskTitle(),
		TeamID:       s.randomTeamID(),
		TeamName:     s.randomTeamName(),
		AssigneeID:   s.randomUserID(),
		AssigneeName: s.randomUserName(),
	}
}
