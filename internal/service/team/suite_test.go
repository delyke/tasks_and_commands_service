package team

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	teamMocks "github.com/delyke/tasks_and_commands_service/internal/domain/team/mocks"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
	userMocks "github.com/delyke/tasks_and_commands_service/internal/domain/user/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx        context.Context //nolint:containedctx
	teamRepo   *teamMocks.Repository
	memberRepo *teamMocks.MemberRepository
	userRepo   *userMocks.Repository
	service    *Service
	faker      *gofakeit.Faker
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.teamRepo = teamMocks.NewRepository(s.T())
	s.memberRepo = teamMocks.NewMemberRepository(s.T())
	s.userRepo = userMocks.NewRepository(s.T())
	s.service = NewTeamService(s.teamRepo, s.memberRepo, s.userRepo)
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

func (s *ServiceSuite) randomMemberID() team.MemberID {
	return team.MemberID(s.faker.Uint64())
}

func (s *ServiceSuite) randomEmail() string {
	return s.faker.Email()
}

func (s *ServiceSuite) randomTeamName() string {
	return s.faker.Company()
}

func (s *ServiceSuite) randomDescription() string {
	return s.faker.Sentence(10)
}

func (s *ServiceSuite) generateUser() *user.User {
	email, _ := user.NewEmail(s.randomEmail())
	return &user.User{
		ID:           s.randomUserID(),
		Email:        email,
		PasswordHash: s.faker.Password(true, true, true, false, false, 32),
		Name:         s.faker.Name(),
	}
}

func (s *ServiceSuite) generateTeam(createdBy user.ID) *team.Team {
	return &team.Team{
		ID:          s.randomTeamID(),
		Name:        s.randomTeamName(),
		Description: s.randomDescription(),
		CreatedBy:   createdBy,
	}
}

func (s *ServiceSuite) generateMember(teamID team.ID, userID user.ID, role team.Role) *team.Member {
	return &team.Member{
		ID:     s.randomMemberID(),
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}
}
