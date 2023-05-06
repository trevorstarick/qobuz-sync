package main

import "github.com/pkg/errors"

var (
	ErrUnavailable    = errors.New("file unavailable")
	ErrAlreadyExists  = errors.New("file already exists")
	ErrAuthFailed     = errors.New("authentication failed")
	ErrInvalidArgs    = errors.New("invalid arguments")
	ErrNotImplemented = errors.New("not implemented")
)
