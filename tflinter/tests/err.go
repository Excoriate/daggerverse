package main

import (
	"errors"
	"fmt"
	"strings"
)

// Constants for module name and emoji representations for different types of messages.
const (
	ModuleName = "TflinterTests🧪"
	ErrorEmoji = "❌"
)

// ModuleTestError represents a custom error for the GoToolbox module.
type ModuleTestError struct {
	message string
	err     error
}

// Error returns the error message with a prefixed module name and emoji.
func (e *ModuleTestError) Error() string {
	prefix := fmt.Sprintf("%s [%s]", ErrorEmoji, ModuleName)
	if e.err != nil {
		return fmt.Sprintf("%s %s: %v", prefix, e.message, e.err)
	}

	return fmt.Sprintf("%s %s", prefix, e.message)
}

// Unwrap returns the underlying error.
func (e *ModuleTestError) Unwrap() error {
	return e.err
}

// NewError creates a new ModuleTestError with a custom message.
//
// Parameters:
//   - message: Custom error message.
//
// Returns:
//   - *ModuleTestError: A new ModuleTestError instance.
func NewError(message string) *ModuleTestError {
	return &ModuleTestError{message: message}
}

// WrapError wraps an existing error with a custom message.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - message: Custom error message.
//
// Returns:
//   - *ModuleTestError: A wrapped error with a custom message.
func WrapError(err error, message string) *ModuleTestError {
	return &ModuleTestError{
		message: message,
		err:     err,
	}
}

// Errorf creates a new ModuleTestError with a formatted message.
//
// Parameters:
//   - format: The format string for the error message.
//   - args: The arguments to be substituted in the format string.
//
// Returns:
//   - *ModuleTestError: A new ModuleTestError instance with a formatted message.
func Errorf(format string, args ...interface{}) *ModuleTestError {
	return &ModuleTestError{
		message: fmt.Sprintf(format, args...),
	}
}

// WrapErrorf wraps an existing error with a formatted message.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - format: The format string for the error message.
//   - args: The arguments to be substituted in the format string.
//
// Returns:
//   - *ModuleTestError: A wrapped error with a formatted message.
func WrapErrorf(err error, format string, args ...interface{}) *ModuleTestError {
	return &ModuleTestError{
		message: fmt.Sprintf(format, args...),
		err:     err,
	}
}

// JoinErrors joins multiple errors into a single ModuleError.
//
// Parameters:
//   - errs: A variadic list of errors to be joined.
//
// Returns:
//   - *ModuleError: A new ModuleError containing the joined error messages,
//     or nil if no errors were provided.
func JoinErrors(errs ...error) *ModuleTestError {
	if len(errs) == 0 {
		return nil
	}

	messages := make([]string, 0, len(errs))

	for _, err := range errs {
		if err != nil {
			var me *ModuleTestError
			if errors.As(err, &me) {
				// If it's already a ModuleError, strip the prefix
				messages = append(messages, strings.TrimPrefix(me.Error(), fmt.Sprintf("%s [%s] ", ErrorEmoji, ModuleName)))
			}
		}
	}

	return &ModuleTestError{
		message: strings.Join(messages, "\n"),
	}
}
