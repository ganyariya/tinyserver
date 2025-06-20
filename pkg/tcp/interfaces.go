package tcp

import (
	"io"
	"net"
	"time"
)

// Connection represents a TCP connection interface
type Connection interface {
	// Read reads data from the connection
	Read([]byte) (int, error)

	// Write writes data to the connection
	Write([]byte) (int, error)

	// Close closes the connection
	Close() error

	// LocalAddr returns the local network address
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address
	RemoteAddr() net.Addr

	// SetDeadline sets the read and write deadlines
	SetDeadline(time.Time) error

	// SetReadDeadline sets the deadline for future Read calls
	SetReadDeadline(time.Time) error

	// SetWriteDeadline sets the deadline for future Write calls
	SetWriteDeadline(time.Time) error
}

// Listener represents a TCP listener interface
type Listener interface {
	// Accept waits for and returns the next connection to the listener
	Accept() (Connection, error)

	// Close closes the listener
	Close() error

	// Addr returns the listener's network address
	Addr() net.Addr
}

// Dialer represents a TCP dialer interface for creating outbound connections
type Dialer interface {
	// Dial connects to the address on the named network
	Dial(network, address string) (Connection, error)

	// DialTimeout acts like Dial but takes a timeout
	DialTimeout(network, address string, timeout time.Duration) (Connection, error)
}

// Server represents a TCP server interface
type Server interface {
	// Start starts the server
	Start() error

	// Stop stops the server
	Stop() error

	// IsRunning returns true if the server is running
	IsRunning() bool

	// Addr returns the server's listening address
	Addr() net.Addr

	// SetHandler sets the connection handler function
	SetHandler(ConnectionHandler)
}

// ConnectionHandler represents a function that handles incoming connections
type ConnectionHandler func(Connection)

// Client represents a TCP client interface
type Client interface {
	// Connect establishes a connection to the server
	Connect(address string) error

	// ConnectWithTimeout establishes a connection with a timeout
	ConnectWithTimeout(address string, timeout time.Duration) error

	// Disconnect closes the connection
	Disconnect() error

	// IsConnected returns true if the client is connected
	IsConnected() bool

	// Send sends data to the server
	Send([]byte) error

	// Receive receives data from the server
	Receive([]byte) (int, error)

	// GetConnection returns the underlying connection
	GetConnection() Connection
}

// ConnectionFactory creates new connections
type ConnectionFactory interface {
	// CreateListener creates a new listener
	CreateListener(network, address string) (Listener, error)

	// CreateDialer creates a new dialer
	CreateDialer() Dialer

	// WrapConnection wraps a net.Conn into our Connection interface
	WrapConnection(net.Conn) Connection
}

// ConnectionPool manages a pool of connections
type ConnectionPool interface {
	// Get retrieves a connection from the pool
	Get() (Connection, error)

	// Put returns a connection to the pool
	Put(Connection) error

	// Close closes all connections in the pool
	Close() error

	// Size returns the current size of the pool
	Size() int

	// Available returns the number of available connections
	Available() int
}

// ConnectionMultiplexer handles multiple connections
type ConnectionMultiplexer interface {
	// AddConnection adds a connection to be multiplexed
	AddConnection(Connection) error

	// RemoveConnection removes a connection from multiplexing
	RemoveConnection(Connection) error

	// Broadcast sends data to all connections
	Broadcast([]byte) error

	// GetConnections returns all active connections
	GetConnections() []Connection

	// GetConnectionCount returns the number of active connections
	GetConnectionCount() int

	// Close closes all connections
	Close() error
}

// MessageReader provides message-based reading from connections
type MessageReader interface {
	// ReadMessage reads a complete message from the connection
	ReadMessage(Connection) ([]byte, error)

	// ReadMessageWithTimeout reads a message with a timeout
	ReadMessageWithTimeout(Connection, time.Duration) ([]byte, error)

	// SetMessageDelimiter sets the delimiter for message boundaries
	SetMessageDelimiter([]byte)
}

// MessageWriter provides message-based writing to connections
type MessageWriter interface {
	// WriteMessage writes a complete message to the connection
	WriteMessage(Connection, []byte) error

	// WriteMessageWithTimeout writes a message with a timeout
	WriteMessageWithTimeout(Connection, []byte, time.Duration) error

	// SetMessageDelimiter sets the delimiter for message boundaries
	SetMessageDelimiter([]byte)
}

// BufferedConnection provides buffered I/O operations
type BufferedConnection interface {
	Connection

	// BufferedReader returns a buffered reader for the connection
	BufferedReader() io.Reader

	// BufferedWriter returns a buffered writer for the connection
	BufferedWriter() io.Writer

	// Flush flushes any buffered data
	Flush() error

	// ReadLine reads a line from the connection
	ReadLine() ([]byte, error)

	// WriteLine writes a line to the connection
	WriteLine([]byte) error
}
