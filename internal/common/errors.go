package common

import (
	"fmt"
)

// ErrorType represents the type of error that occurred
type ErrorType int

const (
	// ErrorTypeNetwork represents network-related errors
	ErrorTypeNetwork ErrorType = iota
	// ErrorTypeProtocol represents protocol parsing errors
	ErrorTypeProtocol
	// ErrorTypeServer represents server-specific errors
	ErrorTypeServer
	// ErrorTypeClient represents client-specific errors
	ErrorTypeClient
	// ErrorTypeIO represents I/O related errors
	ErrorTypeIO
	// ErrorTypeTimeout represents timeout errors
	ErrorTypeTimeout
	// ErrorTypeInvalidInput represents invalid input errors
	ErrorTypeInvalidInput
)

// String returns the string representation of ErrorType
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeNetwork:
		return "NETWORK"
	case ErrorTypeProtocol:
		return "PROTOCOL"
	case ErrorTypeServer:
		return "SERVER"
	case ErrorTypeClient:
		return "CLIENT"
	case ErrorTypeIO:
		return "IO"
	case ErrorTypeTimeout:
		return "TIMEOUT"
	case ErrorTypeInvalidInput:
		return "INVALID_INPUT"
	default:
		return "UNKNOWN"
	}
}

// TinyServerError represents a custom error type for TinyServer
type TinyServerError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// Error implements the error interface
func (e *TinyServerError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *TinyServerError) Unwrap() error {
	return e.Cause
}

// NewError creates a new TinyServerError
func NewError(errorType ErrorType, message string) *TinyServerError {
	return &TinyServerError{
		Type:    errorType,
		Message: message,
	}
}

// NewErrorWithCause creates a new TinyServerError with a cause
func NewErrorWithCause(errorType ErrorType, message string, cause error) *TinyServerError {
	return &TinyServerError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// Common error constructors for frequently used errors

// NetworkError creates a network-related error
func NetworkError(message string) *TinyServerError {
	return NewError(ErrorTypeNetwork, message)
}

// NetworkErrorWithCause creates a network-related error with a cause
func NetworkErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeNetwork, message, cause)
}

// ProtocolError creates a protocol parsing error
func ProtocolError(message string) *TinyServerError {
	return NewError(ErrorTypeProtocol, message)
}

// ProtocolErrorWithCause creates a protocol parsing error with a cause
func ProtocolErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeProtocol, message, cause)
}

// ServerError creates a server-specific error
func ServerError(message string) *TinyServerError {
	return NewError(ErrorTypeServer, message)
}

// ServerErrorWithCause creates a server-specific error with a cause
func ServerErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeServer, message, cause)
}

// ClientError creates a client-specific error
func ClientError(message string) *TinyServerError {
	return NewError(ErrorTypeClient, message)
}

// ClientErrorWithCause creates a client-specific error with a cause
func ClientErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeClient, message, cause)
}

// IOError creates an I/O related error
func IOError(message string) *TinyServerError {
	return NewError(ErrorTypeIO, message)
}

// IOErrorWithCause creates an I/O related error with a cause
func IOErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeIO, message, cause)
}

// TimeoutError creates a timeout error
func TimeoutError(message string) *TinyServerError {
	return NewError(ErrorTypeTimeout, message)
}

// TimeoutErrorWithCause creates a timeout error with a cause
func TimeoutErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeTimeout, message, cause)
}

// InvalidInputError creates an invalid input error
func InvalidInputError(message string) *TinyServerError {
	return NewError(ErrorTypeInvalidInput, message)
}

// InvalidInputErrorWithCause creates an invalid input error with a cause
func InvalidInputErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeInvalidInput, message, cause)
}

// HTTPError creates an HTTP error
func HTTPError(message string) *TinyServerError {
	return NewError(ErrorTypeProtocol, message)
}

// HTTPErrorWithCause creates an HTTP error with underlying cause
func HTTPErrorWithCause(message string, cause error) *TinyServerError {
	return NewErrorWithCause(ErrorTypeProtocol, message, cause)
}
