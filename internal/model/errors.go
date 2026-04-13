package model

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInternalError  = errors.New("internal server error")
)
