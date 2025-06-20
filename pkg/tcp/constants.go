package tcp

import "time"

// Network protocols
const (
	// NetworkTCP represents the TCP network protocol
	NetworkTCP = "tcp"

	// NetworkTCP4 represents TCP over IPv4
	NetworkTCP4 = "tcp4"

	// NetworkTCP6 represents TCP over IPv6
	NetworkTCP6 = "tcp6"
)

// Default ports
const (
	// DefaultEchoPort is the default port for echo server
	DefaultEchoPort = 8080

	// DefaultHTTPPort is the default HTTP port
	DefaultHTTPPort = 80

	// DefaultHTTPSPort is the default HTTPS port
	DefaultHTTPSPort = 443
)

// Connection settings
const (
	// DefaultDialTimeout is the default timeout for dialing connections
	DefaultDialTimeout = 30 * time.Second

	// DefaultKeepAlive is the default keep-alive period
	DefaultKeepAlive = 15 * time.Second

	// DefaultReadBufferSize is the default read buffer size
	DefaultReadBufferSize = 4096

	// DefaultWriteBufferSize is the default write buffer size
	DefaultWriteBufferSize = 4096

	// MaxReadBufferSize is the maximum read buffer size
	MaxReadBufferSize = 65536

	// MaxWriteBufferSize is the maximum write buffer size
	MaxWriteBufferSize = 65536
)

// Connection pool settings
const (
	// DefaultPoolSize is the default connection pool size
	DefaultPoolSize = 10

	// MaxPoolSize is the maximum connection pool size
	MaxPoolSize = 100

	// DefaultPoolTimeout is the default timeout for getting connections from pool
	DefaultPoolTimeout = 5 * time.Second
)

// Message settings
const (
	// DefaultMessageDelimiter is the default message delimiter
	DefaultMessageDelimiter = "\n"

	// MaxMessageSize is the maximum message size in bytes
	MaxMessageSize = 1024 * 1024 // 1MB

	// DefaultMessageTimeout is the default timeout for message operations
	DefaultMessageTimeout = 30 * time.Second
)

// Server settings
const (
	// DefaultMaxConnections is the default maximum number of concurrent connections
	DefaultMaxConnections = 1000

	// DefaultServerReadTimeout is the default server read timeout
	DefaultServerReadTimeout = 30 * time.Second

	// DefaultServerWriteTimeout is the default server write timeout
	DefaultServerWriteTimeout = 30 * time.Second

	// DefaultServerIdleTimeout is the default server idle timeout
	DefaultServerIdleTimeout = 60 * time.Second
)

// Buffer sizes for different operations
const (
	// SmallBufferSize for small operations
	SmallBufferSize = 512

	// MediumBufferSize for medium operations
	MediumBufferSize = 2048

	// LargeBufferSize for large operations
	LargeBufferSize = 8192

	// HugeBufferSize for very large operations
	HugeBufferSize = 32768
)

// Connection states
const (
	// StateDisconnected represents a disconnected state
	StateDisconnected = "disconnected"

	// StateConnecting represents a connecting state
	StateConnecting = "connecting"

	// StateConnected represents a connected state
	StateConnected = "connected"

	// StateClosing represents a closing state
	StateClosing = "closing"

	// StateError represents an error state
	StateError = "error"
)

// Error messages specific to TCP operations
const (
	// ErrMsgConnectionClosed indicates the connection is closed
	ErrMsgConnectionClosed = "connection is closed"

	// ErrMsgConnectionTimeout indicates a connection timeout
	ErrMsgConnectionTimeout = "connection timeout"

	// ErrMsgInvalidAddress indicates an invalid address
	ErrMsgInvalidAddress = "invalid address"

	// ErrMsgListenerClosed indicates the listener is closed
	ErrMsgListenerClosed = "listener is closed"

	// ErrMsgMaxConnectionsReached indicates maximum connections reached
	ErrMsgMaxConnectionsReached = "maximum connections reached"

	// ErrMsgPoolExhausted indicates the connection pool is exhausted
	ErrMsgPoolExhausted = "connection pool exhausted"

	// ErrMsgMessageTooLarge indicates the message is too large
	ErrMsgMessageTooLarge = "message too large"

	// ErrMsgInvalidMessageFormat indicates invalid message format
	ErrMsgInvalidMessageFormat = "invalid message format"
)
