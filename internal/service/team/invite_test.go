package team

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

func (s *ServiceSuite) TestInviteMember_Success() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	invitee := s.generateUser()
	memberID := s.randomMemberID()

	inviterMember := s.generateMember(teamID, inviterID, team.RoleOwner)

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: invitee.Email.String(),
		Role:         "member",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, invitee.Email).
		Return(invitee, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, invitee.ID).
		Return(false, nil)

	s.memberRepo.EXPECT().
		Add(s.ctx, mock.MatchedBy(func(m *team.Member) bool {
			return m.TeamID == teamID &&
				m.UserID == invitee.ID &&
				m.Role == team.RoleMember
		})).
		Return(memberID, nil)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.NoError(err)
}

func (s *ServiceSuite) TestInviteMember_InviterNotFound() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: s.randomEmail(),
		Role:         "member",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(nil, domain.ErrNotFound)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestInviteMember_InviterCannotInvite() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	inviterMember := s.generateMember(teamID, inviterID, team.RoleMember) // участник не может приглашать

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: s.randomEmail(),
		Role:         "member",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, domain.ErrForbidden)
}

func (s *ServiceSuite) TestInviteMember_InvalidEmail() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	inviterMember := s.generateMember(teamID, inviterID, team.RoleAdmin)

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: "invalid-email",
		Role:         "member",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, user.ErrInvalidEmail)
}

func (s *ServiceSuite) TestInviteMember_InviteeNotFound() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	inviterMember := s.generateMember(teamID, inviterID, team.RoleAdmin)
	inviteeEmail, _ := user.NewEmail(s.randomEmail())

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: inviteeEmail.String(),
		Role:         "member",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, inviteeEmail).
		Return(nil, domain.ErrNotFound)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, domain.ErrNotFound)
}

func (s *ServiceSuite) TestInviteMember_AlreadyMember() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	invitee := s.generateUser()
	inviterMember := s.generateMember(teamID, inviterID, team.RoleOwner)

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: invitee.Email.String(),
		Role:         "member",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, invitee.Email).
		Return(invitee, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, invitee.ID).
		Return(true, nil)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.Error(err)
	s.ErrorIs(err, domain.ErrAlreadyExists)
}

func (s *ServiceSuite) TestInviteMember_OwnerRoleDowngradedToAdmin() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	invitee := s.generateUser()
	memberID := s.randomMemberID()
	inviterMember := s.generateMember(teamID, inviterID, team.RoleOwner)

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: invitee.Email.String(),
		Role:         "owner", // пытаемся пригласить как владельца
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, invitee.Email).
		Return(invitee, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, invitee.ID).
		Return(false, nil)

	// Должно понизиться до админа
	s.memberRepo.EXPECT().
		Add(s.ctx, mock.MatchedBy(func(m *team.Member) bool {
			return m.Role == team.RoleAdmin // понижено с owner
		})).
		Return(memberID, nil)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.NoError(err)
}

func (s *ServiceSuite) TestInviteMember_InvalidRoleDefaultsToMember() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()
	invitee := s.generateUser()
	memberID := s.randomMemberID()
	inviterMember := s.generateMember(teamID, inviterID, team.RoleOwner)

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: invitee.Email.String(),
		Role:         "invalid-role",
	}

	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(inviterMember, nil)

	s.userRepo.EXPECT().
		GetByEmail(s.ctx, invitee.Email).
		Return(invitee, nil)

	s.memberRepo.EXPECT().
		IsMember(s.ctx, teamID, invitee.ID).
		Return(false, nil)

	s.memberRepo.EXPECT().
		Add(s.ctx, mock.MatchedBy(func(m *team.Member) bool {
			return m.Role == team.RoleMember // по умолчанию member
		})).
		Return(memberID, nil)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.NoError(err)
}

func (s *ServiceSuite) TestInviteMember_MemberRepoError() {
	// Подготовка
	teamID := s.randomTeamID()
	inviterID := s.randomUserID()

	input := InviteInput{
		TeamID:       teamID,
		InviterID:    inviterID,
		InviteeEmail: s.randomEmail(),
		Role:         "member",
	}

	expectedErr := errors.New("ошибка базы данных")
	s.memberRepo.EXPECT().
		GetByTeamAndUser(s.ctx, teamID, inviterID).
		Return(nil, expectedErr)

	// Действие
	err := s.service.InviteMember(s.ctx, input)

	// Проверка
	s.Error(err)
	s.Equal(expectedErr, err)
}
