package tcp

import (
	"bufio"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ganyariya/tinyserver/internal/common"
	pkgtcp "github.com/ganyariya/tinyserver/pkg/tcp"
)

// tcpConnection implements the tcp.Connection interface
type tcpConnection struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	logger *common.Logger
	mu     sync.RWMutex
	closed bool
}

// NewConnection creates a new TCP connection wrapper
func NewConnection(conn net.Conn) pkgtcp.Connection {
	return &tcpConnection{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, bufferedReaderSize),
		writer: bufio.NewWriterSize(conn, bufferedWriterSize),
		logger: common.NewDefaultLogger(),
	}
}

// Read reads data from the connection
func (c *tcpConnection) Read(p []byte) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, common.NetworkError("connection is closed")
	}

	return c.conn.Read(p)
}

// Write writes data to the connection
func (c *tcpConnection) Write(p []byte) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, common.NetworkError("connection is closed")
	}

	return c.conn.Write(p)
}

// Close closes the connection
func (c *tcpConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	// Flush any remaining buffered data
	if c.writer != nil {
		if err := c.writer.Flush(); err != nil {
			c.logger.Warn("failed to flush writer during close: %v", err)
		}
	}

	return c.conn.Close()
}

// LocalAddr returns the local network address
func (c *tcpConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address
func (c *tcpConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines
func (c *tcpConnection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
func (c *tcpConnection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
func (c *tcpConnection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// bufferedConnection implements the tcp.BufferedConnection interface
type bufferedConnection struct {
	*tcpConnection
}

// NewBufferedConnection creates a new buffered TCP connection
func NewBufferedConnection(conn net.Conn) pkgtcp.BufferedConnection {
	tcpConn := NewConnection(conn).(*tcpConnection)
	return &bufferedConnection{tcpConnection: tcpConn}
}

// BufferedReader returns a buffered reader for the connection
func (c *bufferedConnection) BufferedReader() io.Reader {
	return c.reader
}

// BufferedWriter returns a buffered writer for the connection
func (c *bufferedConnection) BufferedWriter() io.Writer {
	return c.writer
}

// Flush flushes any buffered data
func (c *bufferedConnection) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return common.NetworkError("connection is closed")
	}

	return c.writer.Flush()
}

// ReadLine reads a line from the connection
func (c *bufferedConnection) ReadLine() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, common.NetworkError("connection is closed")
	}

	// Set read timeout
	if err := c.conn.SetReadDeadline(time.Now().Add(lineReadTimeout)); err != nil {
		return nil, common.NetworkErrorWithCause("failed to set read deadline", err)
	}

	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, common.NetworkErrorWithCause("failed to read line", err)
	}

	// Remove trailing newline
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
		// Remove trailing carriage return if present
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
	}

	return line, nil
}

// WriteLine writes a line to the connection
func (c *bufferedConnection) WriteLine(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return common.NetworkError("connection is closed")
	}

	// Set write timeout
	if err := c.conn.SetWriteDeadline(time.Now().Add(flushTimeout)); err != nil {
		return common.NetworkErrorWithCause("failed to set write deadline", err)
	}

	// Write data followed by CRLF
	if _, err := c.writer.Write(data); err != nil {
		return common.NetworkErrorWithCause("failed to write data", err)
	}

	if _, err := c.writer.Write([]byte("\r\n")); err != nil {
		return common.NetworkErrorWithCause("failed to write line ending", err)
	}

	return c.writer.Flush()
}

// messageConnection provides message-based I/O operations
type messageConnection struct {
	pkgtcp.Connection
	delimiter []byte
	logger    *common.Logger
}

// NewMessageConnection creates a new message-based connection
func NewMessageConnection(conn pkgtcp.Connection) *messageConnection {
	return &messageConnection{
		Connection: conn,
		delimiter:  []byte(pkgtcp.DefaultMessageDelimiter),
		logger:     common.NewDefaultLogger(),
	}
}

// ReadMessage reads a complete message from the connection
func (c *messageConnection) ReadMessage() ([]byte, error) {
	return c.ReadMessageWithTimeout(common.DefaultTimeout)
}

// ReadMessageWithTimeout reads a message with a timeout
func (c *messageConnection) ReadMessageWithTimeout(timeout time.Duration) ([]byte, error) {
	// Set read deadline
	if err := c.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, common.NetworkErrorWithCause("failed to set read deadline", err)
	}

	var buffer []byte
	readBuffer := make([]byte, messageReadChunkSize)

	for {
		n, err := c.Read(readBuffer)
		if err != nil {
			if err == io.EOF && len(buffer) > 0 {
				// Return partial message on EOF
				return buffer, nil
			}
			return nil, common.NetworkErrorWithCause("failed to read message chunk", err)
		}

		buffer = append(buffer, readBuffer[:n]...)

		// Check for message delimiter
		if delimiterIndex := findDelimiter(buffer, c.delimiter); delimiterIndex != -1 {
			message := buffer[:delimiterIndex]
			// Note: In a real implementation, we'd need to handle remaining data
			return message, nil
		}

		// Check message size limit
		if len(buffer) > pkgtcp.MaxMessageSize {
			return nil, common.ProtocolError("message too large")
		}
	}
}

// WriteMessage writes a complete message to the connection
func (c *messageConnection) WriteMessage(data []byte) error {
	return c.WriteMessageWithTimeout(data, common.DefaultTimeout)
}

// WriteMessageWithTimeout writes a message with a timeout
func (c *messageConnection) WriteMessageWithTimeout(data []byte, timeout time.Duration) error {
	// Set write deadline
	if err := c.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return common.NetworkErrorWithCause("failed to set write deadline", err)
	}

	// Write message followed by delimiter
	if _, err := c.Write(data); err != nil {
		return common.NetworkErrorWithCause("failed to write message data", err)
	}

	if _, err := c.Write(c.delimiter); err != nil {
		return common.NetworkErrorWithCause("failed to write message delimiter", err)
	}

	return nil
}

// SetMessageDelimiter sets the delimiter for message boundaries
func (c *messageConnection) SetMessageDelimiter(delimiter []byte) {
	c.delimiter = delimiter
}

// Helper functions

// findDelimiter finds the delimiter in the buffer
func findDelimiter(buffer, delimiter []byte) int {
	if len(delimiter) == 0 {
		return -1
	}

	for i := 0; i <= len(buffer)-len(delimiter); i++ {
		if matchDelimiter(buffer[i:], delimiter) {
			return i
		}
	}

	return -1
}

// matchDelimiter checks if the buffer starts with the delimiter
func matchDelimiter(buffer, delimiter []byte) bool {
	if len(buffer) < len(delimiter) {
		return false
	}

	for i := 0; i < len(delimiter); i++ {
		if buffer[i] != delimiter[i] {
			return false
		}
	}

	return true
}

// configureConnection applies optimal TCP settings to a connection
func configureConnection(conn net.Conn) error {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		// Enable TCP_NODELAY to disable Nagle's algorithm
		if err := tcpConn.SetNoDelay(tcpNoDelay); err != nil {
			return common.NetworkErrorWithCause("failed to set TCP_NODELAY", err)
		}

		// Enable keep-alive
		if err := tcpConn.SetKeepAlive(tcpKeepAlive); err != nil {
			return common.NetworkErrorWithCause("failed to set keep-alive", err)
		}

		// Set keep-alive period
		if tcpKeepAlive {
			if err := tcpConn.SetKeepAlivePeriod(tcpKeepAlivePeriod); err != nil {
				return common.NetworkErrorWithCause("failed to set keep-alive period", err)
			}
		}
	}

	return nil
}
