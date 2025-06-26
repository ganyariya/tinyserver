package http

import "time"

// Internal HTTP processing constants
const (
	// DefaultBufferSize is the default buffer size for HTTP processing
	DefaultBufferSize = 4096

	// MaxHeaderLines is the maximum number of header lines
	MaxHeaderLines = 100

	// MaxRequestLineLength is the maximum length of the request line
	MaxRequestLineLength = 2048

	// MaxHeaderLineLength is the maximum length of a header line
	MaxHeaderLineLength = 4096

	// MaxChunkSize is the maximum size of a chunk in chunked encoding
	MaxChunkSize = 1 << 16 // 64KB

	// ParserTimeout is the default timeout for parsing operations
	ParserTimeout = 5 * time.Second

	// WriteTimeout is the default timeout for write operations
	WriteTimeout = 5 * time.Second

	// ReadTimeout is the default timeout for read operations
	ReadTimeout = 5 * time.Second
)

// Parser state constants
const (
	// ParserStateRequestLine indicates parsing request line
	ParserStateRequestLine = iota
	// ParserStateHeaders indicates parsing headers
	ParserStateHeaders
	// ParserStateBody indicates parsing body
	ParserStateBody
	// ParserStateComplete indicates parsing is complete
	ParserStateComplete
)

// Chunked encoding constants
const (
	// ChunkSizeHex indicates chunk size is in hexadecimal
	ChunkSizeHex = "0123456789abcdefABCDEF"
	// ChunkExtensionSeparator separates chunk size and extensions
	ChunkExtensionSeparator = ";"
	// ChunkTrailerStart indicates the start of chunk trailers
	ChunkTrailerStart = "0\r\n"
	// ChunkEnd indicates the end of chunked data
	ChunkEnd = "\r\n"
)

// HTTP parsing patterns
const (
	// HTTPMethodPattern is the pattern for HTTP methods
	HTTPMethodPattern = "^[A-Z]+$"
	// HTTPPathPattern is the pattern for HTTP paths
	HTTPPathPattern = "^[^\\s]+$"
	// HTTPVersionPattern is the pattern for HTTP versions
	HTTPVersionPattern = "^HTTP/[0-9]\\.[0-9]$"
	// HTTPHeaderNamePattern is the pattern for header names
	HTTPHeaderNamePattern = "^[a-zA-Z0-9][a-zA-Z0-9\\-]*$"
)

// Error messages
const (
	// ErrInvalidRequestLine indicates invalid request line
	ErrInvalidRequestLine = "invalid HTTP request line"
	// ErrInvalidMethod indicates invalid HTTP method
	ErrInvalidMethod = "invalid HTTP method"
	// ErrInvalidPath indicates invalid request path
	ErrInvalidPath = "invalid request path"
	// ErrInvalidVersion indicates invalid HTTP version
	ErrInvalidVersion = "invalid HTTP version"
	// ErrInvalidHeader indicates invalid header format
	ErrInvalidHeader = "invalid header format"
	// ErrInvalidStatusCode indicates invalid status code
	ErrInvalidStatusCode = "invalid status code"
	// ErrInvalidContentLength indicates invalid content length
	ErrInvalidContentLength = "invalid content length"
	// ErrRequestTooLarge indicates request is too large
	ErrRequestTooLarge = "request too large"
	// ErrHeaderTooLarge indicates header is too large
	ErrHeaderTooLarge = "header too large"
	// ErrChunkedEncodingInvalid indicates invalid chunked encoding
	ErrChunkedEncodingInvalid = "invalid chunked encoding"
	// ErrUnexpectedEOF indicates unexpected end of input
	ErrUnexpectedEOF = "unexpected end of input"
	// ErrParseTimeout indicates parsing timeout
	ErrParseTimeout = "parsing timeout"
)
