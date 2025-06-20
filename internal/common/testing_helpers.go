package common

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// TestHelper provides utility functions for testing
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new TestHelper instance
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertNoError asserts that an error is nil
func (th *TestHelper) AssertNoError(err error, message ...string) {
	th.t.Helper()
	if err != nil {
		msg := "expected no error"
		if len(message) > 0 {
			msg = message[0]
		}
		th.t.Fatalf("%s, but got: %v", msg, err)
	}
}

// AssertError asserts that an error is not nil
func (th *TestHelper) AssertError(err error, message ...string) {
	th.t.Helper()
	if err == nil {
		msg := "expected an error"
		if len(message) > 0 {
			msg = message[0]
		}
		th.t.Fatalf("%s, but got nil", msg)
	}
}

// AssertEqual asserts that two values are equal
func (th *TestHelper) AssertEqual(expected, actual interface{}, message ...string) {
	th.t.Helper()
	if expected != actual {
		msg := "values are not equal"
		if len(message) > 0 {
			msg = message[0]
		}
		th.t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func (th *TestHelper) AssertNotEqual(expected, actual interface{}, message ...string) {
	th.t.Helper()
	if expected == actual {
		msg := "values should not be equal"
		if len(message) > 0 {
			msg = message[0]
		}
		th.t.Fatalf("%s: both values are %v", msg, expected)
	}
}

// AssertStringContains asserts that a string contains a substring
func (th *TestHelper) AssertStringContains(str, substr string, message ...string) {
	th.t.Helper()
	if !bytes.Contains([]byte(str), []byte(substr)) {
		msg := "string does not contain expected substring"
		if len(message) > 0 {
			msg = message[0]
		}
		th.t.Fatalf("%s: '%s' not found in '%s'", msg, substr, str)
	}
}

// AssertBytesEqual asserts that two byte slices are equal
func (th *TestHelper) AssertBytesEqual(expected, actual []byte, message ...string) {
	th.t.Helper()
	if !bytes.Equal(expected, actual) {
		msg := "byte slices are not equal"
		if len(message) > 0 {
			msg = message[0]
		}
		th.t.Fatalf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// MockReader creates a mock reader from a string
func (th *TestHelper) MockReader(data string) io.Reader {
	return bytes.NewReader([]byte(data))
}

// MockWriter creates a mock writer using a bytes.Buffer
func (th *TestHelper) MockWriter() (*bytes.Buffer, io.Writer) {
	buf := &bytes.Buffer{}
	return buf, buf
}

// GetFreePort returns a free port for testing
func (th *TestHelper) GetFreePort() int {
	th.t.Helper()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		th.t.Fatalf("failed to get free port: %v", err)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

// WaitForPort waits for a port to become available (for server startup)
func (th *TestHelper) WaitForPort(port int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", port)))
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}

	return false
}

// CreateTestServer creates a test TCP server for testing
func (th *TestHelper) CreateTestServer(port int, handler func(net.Conn)) (func(), error) {
	th.t.Helper()

	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", port)))
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})

	go func() {
		defer listener.Close()

		for {
			select {
			case <-done:
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					select {
					case <-done:
						return
					default:
						th.t.Logf("accept error: %v", err)
						continue
					}
				}

				go func(c net.Conn) {
					defer c.Close()
					handler(c)
				}(conn)
			}
		}
	}()

	// Return cleanup function
	cleanup := func() {
		close(done)
		listener.Close()
	}

	return cleanup, nil
}

// Package-level helper functions for convenience

// AssertNoError is a package-level helper for asserting no error
func AssertNoError(t *testing.T, err error, message ...string) {
	NewTestHelper(t).AssertNoError(err, message...)
}

// AssertError is a package-level helper for asserting an error
func AssertError(t *testing.T, err error, message ...string) {
	NewTestHelper(t).AssertError(err, message...)
}

// AssertEqual is a package-level helper for asserting equality
func AssertEqual(t *testing.T, expected, actual interface{}, message ...string) {
	NewTestHelper(t).AssertEqual(expected, actual, message...)
}

// MockReader is a package-level helper for creating mock readers
func MockReader(data string) io.Reader {
	return bytes.NewReader([]byte(data))
}

// MockWriter is a package-level helper for creating mock writers
func MockWriter() (*bytes.Buffer, io.Writer) {
	buf := &bytes.Buffer{}
	return buf, buf
}
