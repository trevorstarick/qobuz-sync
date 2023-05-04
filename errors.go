package main

import "github.com/pkg/errors"

var (
	ErrUnavailable   = errors.New("file unavailable")
	ErrAlreadyExists = errors.New("file already exists")
)
