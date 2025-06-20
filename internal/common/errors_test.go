package common

import (
	"errors"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{ErrorTypeNetwork, "NETWORK"},
		{ErrorTypeProtocol, "PROTOCOL"},
		{ErrorTypeServer, "SERVER"},
		{ErrorTypeClient, "CLIENT"},
		{ErrorTypeIO, "IO"},
		{ErrorTypeTimeout, "TIMEOUT"},
		{ErrorTypeInvalidInput, "INVALID_INPUT"},
		{ErrorType(999), "UNKNOWN"},
	}

	for _, test := range tests {
		result := test.errorType.String()
		if result != test.expected {
			t.Errorf("ErrorType(%d).String() = %s, expected %s", test.errorType, result, test.expected)
		}
	}
}

func TestNewError(t *testing.T) {
	err := NewError(ErrorTypeNetwork, "test message")

	if err.Type != ErrorTypeNetwork {
		t.Errorf("Expected error type %v, got %v", ErrorTypeNetwork, err.Type)
	}

	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	if err.Cause != nil {
		t.Errorf("Expected nil cause, got %v", err.Cause)
	}
}

func TestNewErrorWithCause(t *testing.T) {
	cause := errors.New("original error")
	err := NewErrorWithCause(ErrorTypeNetwork, "test message", cause)

	if err.Type != ErrorTypeNetwork {
		t.Errorf("Expected error type %v, got %v", ErrorTypeNetwork, err.Type)
	}

	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, err.Cause)
	}
}

func TestTinyServerError_Error(t *testing.T) {
	// Test without cause
	err := NewError(ErrorTypeNetwork, "test message")
	expected := "[NETWORK] test message"
	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}

	// Test with cause
	cause := errors.New("original error")
	errWithCause := NewErrorWithCause(ErrorTypeNetwork, "test message", cause)
	expectedWithCause := "[NETWORK] test message: original error"
	if errWithCause.Error() != expectedWithCause {
		t.Errorf("Expected error string '%s', got '%s'", expectedWithCause, errWithCause.Error())
	}
}

func TestTinyServerError_Unwrap(t *testing.T) {
	// Test without cause
	err := NewError(ErrorTypeNetwork, "test message")
	if err.Unwrap() != nil {
		t.Errorf("Expected nil unwrapped error, got %v", err.Unwrap())
	}

	// Test with cause
	cause := errors.New("original error")
	errWithCause := NewErrorWithCause(ErrorTypeNetwork, "test message", cause)
	if errWithCause.Unwrap() != cause {
		t.Errorf("Expected unwrapped error %v, got %v", cause, errWithCause.Unwrap())
	}
}

func TestNetworkError(t *testing.T) {
	err := NetworkError("network failed")

	if err.Type != ErrorTypeNetwork {
		t.Errorf("Expected error type %v, got %v", ErrorTypeNetwork, err.Type)
	}

	if err.Message != "network failed" {
		t.Errorf("Expected message 'network failed', got '%s'", err.Message)
	}
}

func TestNetworkErrorWithCause(t *testing.T) {
	cause := errors.New("socket error")
	err := NetworkErrorWithCause("network failed", cause)

	if err.Type != ErrorTypeNetwork {
		t.Errorf("Expected error type %v, got %v", ErrorTypeNetwork, err.Type)
	}

	if err.Message != "network failed" {
		t.Errorf("Expected message 'network failed', got '%s'", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("Expected cause %v, got %v", cause, err.Cause)
	}
}
