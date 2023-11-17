package app

import (
	"errors"
)

// Errors.
var (
	ErrEmailExist           = errors.New("email exist")
	ErrUsernameExist        = errors.New("username exist")
	ErrNotFound             = errors.New("not found")
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrNotDifferent         = errors.New("the values must be different")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidAuth          = errors.New("invalid auth")
	ErrUserIDAndFileIDExist = errors.New("user_id and file_id exist")
	ErrMaxFiles             = errors.New("post can't save new file")
	ErrAccessDenied         = errors.New("access denied")
	ErrInvalidImageFormat   = errors.New("invalid image format")
)
