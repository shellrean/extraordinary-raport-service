package domain

import "errors"

var (
	ErrNotFound = errors.New("item not found")
	ErrInvalidUser = errors.New("invalid user")
	ErrUnauthorized = errors.New("unauthorized user")
	ErrInvalidToken = errors.New("invalid token")
)