package task

import "github.com/delyke/tasks_and_commands_service/internal/domain/user"

// Filter represents task filtering criteria.
type Filter struct {
	TeamID     *uint64
	Status     *Status
	AssigneeID *user.ID
	Limit      int
	Offset     int
}

// NewFilter creates a new Filter with default pagination.
func NewFilter() *Filter {
	return &Filter{
		Limit:  20,
		Offset: 0,
	}
}

// WithTeamID sets the team ID filter.
func (f *Filter) WithTeamID(teamID uint64) *Filter {
	f.TeamID = &teamID
	return f
}

// WithStatus sets the status filter.
func (f *Filter) WithStatus(status Status) *Filter {
	f.Status = &status
	return f
}

// WithAssignee sets the assignee filter.
func (f *Filter) WithAssignee(assigneeID user.ID) *Filter {
	f.AssigneeID = &assigneeID
	return f
}

// WithPagination sets pagination parameters.
func (f *Filter) WithPagination(limit, offset int) *Filter {
	if limit > 0 && limit <= 100 {
		f.Limit = limit
	}
	if offset >= 0 {
		f.Offset = offset
	}
	return f
}
