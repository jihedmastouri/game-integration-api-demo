package shared

type errorCode string

const (
	ValidationError    errorCode = "REQUEST_VALIDATION_ERROR"
	ServiceUnAvailable errorCode = "SERVICE_UNAVAILABLE"
)
