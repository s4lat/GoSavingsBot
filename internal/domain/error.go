package domain

import "errors"

var (
	ErrUnknown       = errors.New("unknown error")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
