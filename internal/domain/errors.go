package domain

import "errors"

// Common domain errors
var (
	ErrNotFound           = errors.New("entity not found")
	ErrAlreadyExists      = errors.New("entity already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
