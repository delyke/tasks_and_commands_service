package team

import "errors"

var ErrInvalidRole = errors.New("invalid team role")

// Role represents a team member role.
type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// NewRole creates and validates a Role value object.
func NewRole(value string) (Role, error) {
	switch Role(value) {
	case RoleOwner, RoleAdmin, RoleMember:
		return Role(value), nil
	default:
		return "", ErrInvalidRole
	}
}

// String returns the string representation of the role.
func (r Role) String() string {
	return string(r)
}

// CanInvite returns true if this role can invite new members.
func (r Role) CanInvite() bool {
	return r == RoleOwner || r == RoleAdmin
}

// CanManageTasks returns true if this role can manage tasks.
func (r Role) CanManageTasks() bool {
	return r == RoleOwner || r == RoleAdmin || r == RoleMember
}

// CanDeleteTeam returns true if this role can delete the team.
func (r Role) CanDeleteTeam() bool {
	return r == RoleOwner
}
