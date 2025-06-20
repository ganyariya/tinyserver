package common

import "time"

// Network constants
const (
	// DefaultBufferSize is the default buffer size for network operations
	DefaultBufferSize = 4096

	// MaxBufferSize is the maximum buffer size for network operations
	MaxBufferSize = 65536

	// DefaultTimeout is the default timeout for network operations
	DefaultTimeout = 30 * time.Second

	// DefaultReadTimeout is the default read timeout
	DefaultReadTimeout = 10 * time.Second

	// DefaultWriteTimeout is the default write timeout
	DefaultWriteTimeout = 10 * time.Second

	// DefaultServerPort is the default port for TinyServer
	DefaultServerPort = 8080

	// DefaultServerHost is the default host for TinyServer
	DefaultServerHost = "localhost"
)

// Protocol constants
const (
	// ProtocolTCP represents the TCP protocol
	ProtocolTCP = "tcp"

	// ProtocolHTTP represents the HTTP protocol
	ProtocolHTTP = "http"

	// HTTPVersion11 represents HTTP/1.1 version
	HTTPVersion11 = "HTTP/1.1"

	// HTTPVersion10 represents HTTP/1.0 version
	HTTPVersion10 = "HTTP/1.0"
)

// Application constants
const (
	// ApplicationName is the name of the application
	ApplicationName = "TinyServer"

	// ApplicationVersion is the version of the application
	ApplicationVersion = "1.0.0"

	// UserAgent is the default user agent string
	UserAgent = ApplicationName + "/" + ApplicationVersion
)

// File system constants
const (
	// DefaultStaticDir is the default directory for static files
	DefaultStaticDir = "./static"

	// DefaultIndexFile is the default index file name
	DefaultIndexFile = "index.html"

	// DefaultFilePermissions is the default file permissions for created files
	DefaultFilePermissions = 0644

	// DefaultDirPermissions is the default directory permissions for created directories
	DefaultDirPermissions = 0755
)

// Performance constants
const (
	// DefaultMaxConnections is the default maximum number of concurrent connections
	DefaultMaxConnections = 100

	// DefaultMaxRequestSize is the default maximum request size in bytes
	DefaultMaxRequestSize = 1024 * 1024 // 1MB

	// DefaultMaxHeaderSize is the default maximum header size in bytes
	DefaultMaxHeaderSize = 8192 // 8KB

	// DefaultKeepAliveTimeout is the default keep-alive timeout
	DefaultKeepAliveTimeout = 60 * time.Second
)

// Error messages
const (
	// ErrMsgInvalidInput represents an invalid input error message
	ErrMsgInvalidInput = "invalid input provided"

	// ErrMsgNetworkFailure represents a network failure error message
	ErrMsgNetworkFailure = "network operation failed"

	// ErrMsgProtocolError represents a protocol error message
	ErrMsgProtocolError = "protocol parsing error"

	// ErrMsgServerError represents a server error message
	ErrMsgServerError = "server operation failed"

	// ErrMsgClientError represents a client error message
	ErrMsgClientError = "client operation failed"

	// ErrMsgTimeout represents a timeout error message
	ErrMsgTimeout = "operation timed out"

	// ErrMsgIOFailure represents an I/O failure error message
	ErrMsgIOFailure = "I/O operation failed"
)

// MIME types
const (
	// MIMETextPlain represents plain text MIME type
	MIMETextPlain = "text/plain"

	// MIMETextHTML represents HTML MIME type
	MIMETextHTML = "text/html"

	// MIMEApplicationJSON represents JSON MIME type
	MIMEApplicationJSON = "application/json"

	// MIMEApplicationOctetStream represents binary MIME type
	MIMEApplicationOctetStream = "application/octet-stream"
)

// Character sets and encodings
const (
	// CharsetUTF8 represents UTF-8 character set
	CharsetUTF8 = "utf-8"

	// EncodingGzip represents gzip encoding
	EncodingGzip = "gzip"

	// EncodingDeflate represents deflate encoding
	EncodingDeflate = "deflate"
)

// Line endings and separators
const (
	// CRLF represents the HTTP line ending
	CRLF = "\r\n"

	// LF represents the Unix line ending
	LF = "\n"

	// HeaderSeparator represents the separator between headers and body
	HeaderSeparator = CRLF + CRLF

	// SpaceSeparator represents a space character
	SpaceSeparator = " "

	// ColonSeparator represents a colon character
	ColonSeparator = ":"
)
