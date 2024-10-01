package main

import (
	"errors"
	"fmt"
	"strings"
)

// Constants for module name and emoji representations for different types of messages.
const (
	ModuleName = "ModuleTemplateLight"
	ErrorEmoji = "❌"
)

// ModuleError represents a custom error for the GoToolbox module.
type ModuleError struct {
	message string
	err     error
}

// Error returns the error message with a prefixed module name and emoji.
func (e *ModuleError) Error() string {
	prefix := fmt.Sprintf("%s [%s]", ErrorEmoji, ModuleName)
	if e.err != nil {
		return fmt.Sprintf("%s %s: %v", prefix, e.message, e.err)
	}

	return fmt.Sprintf("%s %s", prefix, e.message)
}

// Unwrap returns the underlying error.
func (e *ModuleError) Unwrap() error {
	return e.err
}

// NewError creates a new ModuleError with a custom message.
//
// Parameters:
//   - message: Custom error message.
//
// Returns:
//   - *ModuleError: A new ModuleError instance.
func NewError(message string) *ModuleError {
	return &ModuleError{message: message}
}

// WrapError wraps an existing error with a custom message.
//
// Parameters:
//   - err: The original error to be wrapped.
//   - message: Custom error message.
//
// Returns:
//   - *ModuleError: A wrapped error with a custom message.
func WrapError(err error, message string) *ModuleError {
	return &ModuleError{
		message: message,
		err:     err,
	}
}

// Errorf creates a new ModuleError with a formatted message.
//
// Parameters:
//   - format: The format string for the error message.
//   - args: The arguments to be substituted in the format string.
//
// Returns:
//   - *ModuleError: A new ModuleError instance with a formatted message.
func Errorf(format string, args ...interface{}) *ModuleError {
	return &ModuleError{
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
//   - *ModuleError: A wrapped error with a formatted message.
func WrapErrorf(err error, format string, args ...interface{}) *ModuleError {
	return &ModuleError{
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
func JoinErrors(errs ...error) *ModuleError {
	if len(errs) == 0 {
		return nil
	}

	messages := make([]string, 0, len(errs))

	for _, err := range errs {
		if err != nil {
			var me *ModuleError
			if errors.As(err, &me) {
				// If it's already a ModuleError, strip the prefix
				messages = append(messages, strings.TrimPrefix(me.Error(), fmt.Sprintf("%s [%s] ", ErrorEmoji, ModuleName)))
			}
		}
	}

	return &ModuleError{
		message: strings.Join(messages, "\n"),
	}
}
