package common

import "github.com/pkg/errors"

var (
	ErrAlreadyExists  = errors.New("file already exists")
	ErrUnavailable    = errors.New("file unavailable")
	ErrAuthFailed     = errors.New("authentication failed")
	ErrInvalidArgs    = errors.New("invalid arguments")
	ErrNotImplemented = errors.New("not implemented")
)
