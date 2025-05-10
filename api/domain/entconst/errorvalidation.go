package entconst

import (
	"fmt"
)

// ValidationError is a custom error type
type ValidationError struct {
	msg string
}

// NewValidationError a new ValidationError with a custom message
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg: msg}
}

// Error returns the error message for ValidationError
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", e.msg)
}
