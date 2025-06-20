package tcp

import "time"

// Internal TCP implementation constants

// Connection implementation settings
const (
	// connectionReadBufferSize is the internal read buffer size
	connectionReadBufferSize = 4096

	// connectionWriteBufferSize is the internal write buffer size
	connectionWriteBufferSize = 4096

	// connectionCloseTimeout is the timeout for closing connections gracefully
	connectionCloseTimeout = 5 * time.Second

	// connectionHealthCheckInterval is the interval for connection health checks
	connectionHealthCheckInterval = 30 * time.Second
)

// Listener implementation settings
const (
	// listenerBacklog is the maximum number of pending connections
	listenerBacklog = 128

	// listenerAcceptTimeout is the timeout for accepting connections
	listenerAcceptTimeout = 1 * time.Second

	// listenerShutdownTimeout is the timeout for listener shutdown
	listenerShutdownTimeout = 5 * time.Second
)

// Server implementation settings
const (
	// serverStartupTimeout is the timeout for server startup
	serverStartupTimeout = 10 * time.Second

	// serverShutdownTimeout is the timeout for server shutdown
	serverShutdownTimeout = 30 * time.Second

	// serverWorkerPoolSize is the size of the worker pool for handling connections
	serverWorkerPoolSize = 100

	// serverConnectionQueueSize is the size of the connection queue
	serverConnectionQueueSize = 1000
)

// Client implementation settings
const (
	// clientConnectRetries is the number of connection retries
	clientConnectRetries = 3

	// clientRetryDelay is the delay between connection retries
	clientRetryDelay = 1 * time.Second

	// clientHeartbeatInterval is the interval for sending heartbeats
	clientHeartbeatInterval = 30 * time.Second

	// clientReconnectDelay is the delay before attempting to reconnect
	clientReconnectDelay = 5 * time.Second
)

// Performance tuning constants
const (
	// tcpNoDelay controls the TCP_NODELAY socket option
	tcpNoDelay = true

	// tcpKeepAlive controls the SO_KEEPALIVE socket option
	tcpKeepAlive = true

	// tcpKeepAlivePeriod is the keep-alive period
	tcpKeepAlivePeriod = 15 * time.Second

	// tcpLinger controls the SO_LINGER socket option (-1 to disable)
	tcpLinger = -1

	// tcpReceiveBufferSize is the SO_RCVBUF socket option
	tcpReceiveBufferSize = 65536

	// tcpSendBufferSize is the SO_SNDBUF socket option
	tcpSendBufferSize = 65536
)

// Buffered connection settings
const (
	// bufferedReaderSize is the size of the buffered reader
	bufferedReaderSize = 8192

	// bufferedWriterSize is the size of the buffered writer
	bufferedWriterSize = 8192

	// lineReadTimeout is the timeout for reading a line
	lineReadTimeout = 10 * time.Second

	// flushTimeout is the timeout for flushing buffered data
	flushTimeout = 5 * time.Second
)

// Message handling constants
const (
	// messageHeaderSize is the size of the message header
	messageHeaderSize = 4

	// messageReadChunkSize is the chunk size for reading messages
	messageReadChunkSize = 1024

	// messageWriteChunkSize is the chunk size for writing messages
	messageWriteChunkSize = 1024

	// messageScanBufferSize is the buffer size for message scanning
	messageScanBufferSize = 64 * 1024
)

// Error handling constants
const (
	// maxRetryAttempts is the maximum number of retry attempts
	maxRetryAttempts = 3

	// retryBackoffMultiplier is the multiplier for exponential backoff
	retryBackoffMultiplier = 2

	// maxRetryDelay is the maximum delay between retries
	maxRetryDelay = 30 * time.Second

	// errorLogThreshold is the threshold for logging errors
	errorLogThreshold = 5
)

// Connection pool implementation constants
const (
	// poolCleanupInterval is the interval for cleaning up expired connections
	poolCleanupInterval = 1 * time.Minute

	// poolConnectionMaxIdleTime is the maximum idle time for pooled connections
	poolConnectionMaxIdleTime = 10 * time.Minute

	// poolConnectionMaxLifetime is the maximum lifetime for pooled connections
	poolConnectionMaxLifetime = 1 * time.Hour

	// poolGrowthFactor is the factor by which the pool grows
	poolGrowthFactor = 2
)

// Multiplexer implementation constants
const (
	// multiplexerChannelBufferSize is the buffer size for multiplexer channels
	multiplexerChannelBufferSize = 100

	// multiplexerBroadcastTimeout is the timeout for broadcast operations
	multiplexerBroadcastTimeout = 10 * time.Second

	// multiplexerCleanupInterval is the interval for cleaning up dead connections
	multiplexerCleanupInterval = 30 * time.Second
)

// Internal state constants
const (
	// stateIdle represents an idle connection state
	stateIdle = "idle"

	// stateActive represents an active connection state
	stateActive = "active"

	// stateDraining represents a draining connection state
	stateDraining = "draining"

	// stateClosed represents a closed connection state
	stateClosed = "closed"
)

// Logging constants
const (
	// logConnectionEvents controls whether to log connection events
	logConnectionEvents = true

	// logPerformanceMetrics controls whether to log performance metrics
	logPerformanceMetrics = false

	// logDebugInfo controls whether to log debug information
	logDebugInfo = false

	// logSlowOperationThreshold is the threshold for logging slow operations
	logSlowOperationThreshold = 100 * time.Millisecond
)
