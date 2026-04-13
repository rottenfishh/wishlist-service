package model

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNotFound       = errors.New("not found")
	ErrForbidden      = errors.New("forbidden")
	ErrInternalError  = errors.New("internal server error")
	ErrAlreadyBooked  = errors.New("gift already booked")
	ErrNotUpdated     = errors.New("not updated")
)
