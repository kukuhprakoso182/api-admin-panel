package domain

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrForbidden     = errors.New("access forbidden")

	// ErrInvalidReference: record merujuk ke data induk yang tidak ada (FK 1452).
	ErrInvalidReference = errors.New("referenced record does not exist")
	// ErrInUse: record tidak bisa dihapus karena masih dipakai data lain (FK 1451).
	ErrInUse = errors.New("record is still in use")
)
