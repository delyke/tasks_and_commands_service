package user

import "context"

// Repository defines the interface for user persistence.
type Repository interface {
	Create(ctx context.Context, user *User) (ID, error)
	GetByID(ctx context.Context, id ID) (*User, error)
	GetByEmail(ctx context.Context, email Email) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id ID) error
}
