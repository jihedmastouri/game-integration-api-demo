package internal

type ErrorCode string

const (
	Unauthorized    ErrorCode = "RequestValidationError"
	ValidationError ErrorCode = "RequestValidationError"
)
