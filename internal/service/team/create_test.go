package team

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
)

func (s *ServiceSuite) TestCreateTeam_Success() {
	// Подготовка
	userID := s.randomUserID()
	teamID := s.randomTeamID()
	memberID := s.randomMemberID()

	input := CreateTeamInput{
		Name:        s.randomTeamName(),
		Description: s.randomDescription(),
		CreatedBy:   userID,
	}

	s.teamRepo.EXPECT().
		Create(s.ctx, mock.MatchedBy(func(t *team.Team) bool {
			return t.Name == input.Name &&
				t.Description == input.Description &&
				t.CreatedBy == input.CreatedBy
		})).
		Return(teamID, nil)

	s.memberRepo.EXPECT().
		Add(s.ctx, mock.MatchedBy(func(m *team.Member) bool {
			return m.TeamID == teamID &&
				m.UserID == userID &&
				m.Role == team.RoleOwner
		})).
		Return(memberID, nil)

	// Действие
	result, err := s.service.CreateTeam(s.ctx, input)

	// Проверка
	s.NoError(err)
	s.NotNil(result)
	s.Equal(teamID, result.ID)
	s.Equal(input.Name, result.Name)
	s.Equal(input.Description, result.Description)
	s.Equal(userID, result.CreatedBy)
}

func (s *ServiceSuite) TestCreateTeam_TeamRepoError() {
	// Подготовка
	input := CreateTeamInput{
		Name:        s.randomTeamName(),
		Description: s.randomDescription(),
		CreatedBy:   s.randomUserID(),
	}

	expectedErr := errors.New("ошибка базы данных")
	s.teamRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(team.ID(0), expectedErr)

	// Действие
	result, err := s.service.CreateTeam(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}

func (s *ServiceSuite) TestCreateTeam_MemberRepoError() {
	// Подготовка
	userID := s.randomUserID()
	teamID := s.randomTeamID()

	input := CreateTeamInput{
		Name:        s.randomTeamName(),
		Description: s.randomDescription(),
		CreatedBy:   userID,
	}

	s.teamRepo.EXPECT().
		Create(s.ctx, mock.Anything).
		Return(teamID, nil)

	expectedErr := errors.New("ошибка добавления участника")
	s.memberRepo.EXPECT().
		Add(s.ctx, mock.Anything).
		Return(team.MemberID(0), expectedErr)

	// Действие
	result, err := s.service.CreateTeam(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
}
