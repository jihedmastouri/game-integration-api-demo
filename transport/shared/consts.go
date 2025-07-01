package shared

import (
	"errors"
)

type errorCode string

const (
	ValidationError     errorCode = "REQUEST_VALIDATION_ERROR"
	ServiceUnAvailable  errorCode = "SERVICE_UNAVAILABLE"
	InternalServerError errorCode = "INTERNAL_SERVER_ERROR"
	Unauthorized        errorCode = "UNAUTHORIZED"
)

var (
	ErrServiceUnAvailable error = errors.New("SERVICE_UNAVAILABLE")
)
