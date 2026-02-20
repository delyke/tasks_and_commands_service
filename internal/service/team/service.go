package team

import (
	"github.com/delyke/tasks_and_commands_service/internal/domain/team"
	"github.com/delyke/tasks_and_commands_service/internal/domain/user"
)

// Service handles team operations.
type Service struct {
	teamRepo   team.Repository
	memberRepo team.MemberRepository
	userRepo   user.Repository
}

// Ensure Service implements TeamService.
var _ TeamService = (*Service)(nil)

// NewTeamService creates a new TeamService.
func NewTeamService(teamRepo team.Repository, memberRepo team.MemberRepository, userRepo user.Repository) *Service {
	return &Service{
		teamRepo:   teamRepo,
		memberRepo: memberRepo,
		userRepo:   userRepo,
	}
}
