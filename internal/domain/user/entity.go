package user

import (
	"time"
)

// ID represents a user identifier.
type ID uint64

// User represents a domain user entity.
type User struct {
	ID           ID
	Email        Email
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// New creates a new User entity.
func New(email Email, passwordHash, name string) *User {
	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Name:         name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// WithID returns a copy of the user with the given ID.
func (u *User) WithID(id ID) *User {
	u.ID = id
	return u
}
