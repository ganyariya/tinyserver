package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ganyariya/tinyserver/internal/common"
	pkgtcp "github.com/ganyariya/tinyserver/pkg/tcp"
)

// tcpListener implements the tcp.Listener interface
type tcpListener struct {
	listener   net.Listener
	logger     *common.Logger
	mu         sync.RWMutex
	closed     int32 // atomic
	closeChan  chan struct{}
	acceptChan chan acceptResult
}

// acceptResult represents the result of an accept operation
type acceptResult struct {
	conn pkgtcp.Connection
	err  error
}

// NewListener creates a new TCP listener
func NewListener(network, address string) (pkgtcp.Listener, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, common.NetworkErrorWithCause("failed to create listener", err)
	}

	tcpListener := &tcpListener{
		listener:   listener,
		logger:     common.NewDefaultLogger(),
		closeChan:  make(chan struct{}),
		acceptChan: make(chan acceptResult, 1),
	}

	// Start the accept goroutine
	go tcpListener.acceptLoop()

	return tcpListener, nil
}

// Accept waits for and returns the next connection to the listener
func (l *tcpListener) Accept() (pkgtcp.Connection, error) {
	if atomic.LoadInt32(&l.closed) == 1 {
		return nil, common.NetworkError("listener is closed")
	}

	select {
	case result := <-l.acceptChan:
		return result.conn, result.err
	case <-l.closeChan:
		return nil, common.NetworkError("listener is closed")
	}
}

// Close closes the listener
func (l *tcpListener) Close() error {
	// Set closed flag atomically
	if !atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
		return nil // Already closed
	}

	l.logger.Info("Closing TCP listener on %s", l.listener.Addr())

	// Close the close channel to signal shutdown
	close(l.closeChan)

	// Close the underlying listener
	return l.listener.Close()
}

// Addr returns the listener's network address
func (l *tcpListener) Addr() net.Addr {
	return l.listener.Addr()
}

// acceptLoop runs in a separate goroutine to handle accept operations
func (l *tcpListener) acceptLoop() {
	for {
		// Check if we're closed
		if atomic.LoadInt32(&l.closed) == 1 {
			return
		}

		// Set accept timeout to allow periodic checks
		if tcpListener, ok := l.listener.(*net.TCPListener); ok {
			tcpListener.SetDeadline(time.Now().Add(listenerAcceptTimeout))
		}

		conn, err := l.listener.Accept()
		if err != nil {
			// Check if this is a timeout error and we're not closed
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				if atomic.LoadInt32(&l.closed) == 0 {
					continue // Continue accepting
				}
			}

			// Check if we're closing
			select {
			case <-l.closeChan:
				return
			default:
			}

			// Send the error
			select {
			case l.acceptChan <- acceptResult{nil, common.NetworkErrorWithCause("accept failed", err)}:
			case <-l.closeChan:
				return
			}
			continue
		}

		// Configure the connection for optimal performance
		if err := configureConnection(conn); err != nil {
			l.logger.Warn("Failed to configure connection: %v", err)
		}

		// Wrap the connection
		tcpConn := NewConnection(conn)

		l.logger.Debug("Accepted connection from %s", conn.RemoteAddr())

		// Send the connection
		select {
		case l.acceptChan <- acceptResult{tcpConn, nil}:
		case <-l.closeChan:
			conn.Close()
			return
		}
	}
}

// connectionFactory implements the tcp.ConnectionFactory interface
type connectionFactory struct {
	logger *common.Logger
}

// NewConnectionFactory creates a new connection factory
func NewConnectionFactory() pkgtcp.ConnectionFactory {
	return &connectionFactory{
		logger: common.NewDefaultLogger(),
	}
}

// CreateListener creates a new listener
func (f *connectionFactory) CreateListener(network, address string) (pkgtcp.Listener, error) {
	return NewListener(network, address)
}

// CreateDialer creates a new dialer
func (f *connectionFactory) CreateDialer() pkgtcp.Dialer {
	return NewDialer()
}

// WrapConnection wraps a net.Conn into our Connection interface
func (f *connectionFactory) WrapConnection(conn net.Conn) pkgtcp.Connection {
	return NewConnection(conn)
}

// tcpDialer implements the tcp.Dialer interface
type tcpDialer struct {
	dialer *net.Dialer
	logger *common.Logger
}

// NewDialer creates a new TCP dialer
func NewDialer() pkgtcp.Dialer {
	return &tcpDialer{
		dialer: &net.Dialer{
			Timeout:   pkgtcp.DefaultDialTimeout,
			KeepAlive: pkgtcp.DefaultKeepAlive,
		},
		logger: common.NewDefaultLogger(),
	}
}

// Dial connects to the address on the named network
func (d *tcpDialer) Dial(network, address string) (pkgtcp.Connection, error) {
	conn, err := d.dialer.Dial(network, address)
	if err != nil {
		return nil, common.NetworkErrorWithCause("dial failed", err)
	}

	// Configure the connection for optimal performance
	if err := configureConnection(conn); err != nil {
		d.logger.Warn("Failed to configure connection: %v", err)
	}

	d.logger.Debug("Connected to %s", address)

	return NewConnection(conn), nil
}

// DialTimeout acts like Dial but takes a timeout
func (d *tcpDialer) DialTimeout(network, address string, timeout time.Duration) (pkgtcp.Connection, error) {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: pkgtcp.DefaultKeepAlive,
	}

	conn, err := dialer.Dial(network, address)
	if err != nil {
		return nil, common.NetworkErrorWithCause("dial with timeout failed", err)
	}

	// Configure the connection for optimal performance
	if err := configureConnection(conn); err != nil {
		d.logger.Warn("Failed to configure connection: %v", err)
	}

	d.logger.Debug("Connected to %s with timeout %v", address, timeout)

	return NewConnection(conn), nil
}

// tcpServer implements the tcp.Server interface
type tcpServer struct {
	listener pkgtcp.Listener
	handler  pkgtcp.ConnectionHandler
	logger   *common.Logger
	mu       sync.RWMutex
	running  bool
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewServer creates a new TCP server
func NewServer(network, address string) (pkgtcp.Server, error) {
	listener, err := NewListener(network, address)
	if err != nil {
		return nil, err
	}

	return &tcpServer{
		listener: listener,
		logger:   common.NewDefaultLogger(),
		stopChan: make(chan struct{}),
	}, nil
}

// Start starts the server
func (s *tcpServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return common.ServerError("server is already running")
	}

	if s.handler == nil {
		return common.ServerError("no connection handler set")
	}

	s.running = true
	s.logger.Info("Starting TCP server on %s", s.listener.Addr())

	// Start accepting connections
	s.wg.Add(1)
	go s.acceptLoop()

	return nil
}

// Stop stops the server
func (s *tcpServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.logger.Info("Stopping TCP server")
	s.running = false

	// Signal stop
	close(s.stopChan)

	// Close the listener
	if err := s.listener.Close(); err != nil {
		s.logger.Warn("Error closing listener: %v", err)
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Wait with timeout
	select {
	case <-done:
		s.logger.Info("TCP server stopped successfully")
	case <-time.After(serverShutdownTimeout):
		s.logger.Warn("TCP server shutdown timeout")
	}

	return nil
}

// IsRunning returns true if the server is running
func (s *tcpServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// Addr returns the server's listening address
func (s *tcpServer) Addr() net.Addr {
	return s.listener.Addr()
}

// SetHandler sets the connection handler function
func (s *tcpServer) SetHandler(handler pkgtcp.ConnectionHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

// acceptLoop accepts incoming connections and handles them
func (s *tcpServer) acceptLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stopChan:
				return
			default:
				s.logger.Error("Accept error: %v", err)
				continue
			}
		}

		// Handle connection in a separate goroutine
		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection handles a single connection
func (s *tcpServer) handleConnection(conn pkgtcp.Connection) {
	defer s.wg.Done()
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	s.logger.Info("Handling connection from %s", remoteAddr)

	// Call the handler
	s.handler(conn)

	s.logger.Info("Connection from %s closed", remoteAddr)
}
