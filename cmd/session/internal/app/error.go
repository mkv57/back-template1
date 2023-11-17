package app

import (
	"errors"
)

// Errors.
var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrDuplicate       = errors.New("duplicate")
)
