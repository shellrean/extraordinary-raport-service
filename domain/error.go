package domain

import "errors"

var (
	ErrNotFound = errors.New("item not found")
	ErrInvalidUser = errors.New("invalid user")
)