package usecase

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrForbidden     = errors.New("access forbidden")
)
