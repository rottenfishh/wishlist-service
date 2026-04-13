package model

import "errors"

var (
	InvalidRequest = errors.New("invalid request")
	Unauthorized   = errors.New("unauthorized")
	NotFound       = errors.New("not found")
	Forbidden      = errors.New("forbidden")
	InternalError  = errors.New("internal server error")
)
