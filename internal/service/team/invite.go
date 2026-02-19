package team

import (
	"context"
	"errors"
	"github.com/delyke/tasks_and_commands_service/internal/domain"
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// InviteInput represents team invitation input.
type InviteInput struct {
	TeamID       team.ID
	InviterID    user.ID
	InviteeEmail string
	Role         string
}

// InviteMember invites a user to a team.
func (s *Service) InviteMember(ctx context.Context, input InviteInput) error {
	// Check inviter permissions
	inviterMember, err := s.memberRepo.GetByTeamAndUser(ctx, input.TeamID, input.InviterID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrForbidden
		}
		return err
	}

	if !inviterMember.Role.CanInvite() {
		return domain.ErrForbidden
	}

	// Find invitee by email
	email, err := user.NewEmail(input.InviteeEmail)
	if err != nil {
		return err
	}

	invitee, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Check if already a member
	isMember, err := s.memberRepo.IsMember(ctx, input.TeamID, invitee.ID)
	if err != nil {
		return err
	}
	if isMember {
		return domain.ErrAlreadyExists
	}

	// Parse role
	role, err := team.NewRole(input.Role)
	if err != nil {
		role = team.RoleMember
	}

	// Prevent inviting as owner
	if role == team.RoleOwner {
		role = team.RoleAdmin
	}

	// Add member
	member := team.NewMember(input.TeamID, invitee.ID, role)
	_, err = s.memberRepo.Add(ctx, member)
	return err
}
