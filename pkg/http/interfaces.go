package http

import (
	"io"
	"net"
	"time"
)

// Method represents HTTP methods
type Method string

const (
	// MethodGet represents HTTP GET method
	MethodGet Method = "GET"
	// MethodPost represents HTTP POST method
	MethodPost Method = "POST"
	// MethodPut represents HTTP PUT method
	MethodPut Method = "PUT"
	// MethodDelete represents HTTP DELETE method
	MethodDelete Method = "DELETE"
	// MethodHead represents HTTP HEAD method
	MethodHead Method = "HEAD"
	// MethodOptions represents HTTP OPTIONS method
	MethodOptions Method = "OPTIONS"
	// MethodPatch represents HTTP PATCH method
	MethodPatch Method = "PATCH"
)

// Version represents HTTP version
type Version string

const (
	// Version10 represents HTTP/1.0
	Version10 Version = "HTTP/1.0"
	// Version11 represents HTTP/1.1
	Version11 Version = "HTTP/1.1"
)

// StatusCode represents HTTP status codes
type StatusCode int

// Header represents HTTP headers as key-value pairs
type Header map[string][]string

// Request represents an HTTP request
type Request interface {
	// Method returns the HTTP method
	Method() Method

	// Path returns the request path
	Path() string

	// Version returns the HTTP version
	Version() Version

	// Headers returns the request headers
	Headers() Header

	// Body returns the request body reader
	Body() io.Reader

	// QueryParams returns query parameters
	QueryParams() map[string]string

	// SetMethod sets the HTTP method
	SetMethod(Method)

	// SetPath sets the request path
	SetPath(string)

	// SetVersion sets the HTTP version
	SetVersion(Version)

	// SetHeader sets a header value
	SetHeader(string, string)

	// AddHeader adds a header value
	AddHeader(string, string)

	// SetBody sets the request body
	SetBody(io.Reader)

	// ContentLength returns the content length
	ContentLength() int64

	// RemoteAddr returns the remote address
	RemoteAddr() net.Addr

	// GetHeader returns the first value of the header
	GetHeader(string) string

	// GetHeaders returns all values for the header
	GetHeaders(string) []string

	// HasHeader checks if a header exists
	HasHeader(string) bool
}

// Response represents an HTTP response
type Response interface {
	// StatusCode returns the HTTP status code
	StatusCode() StatusCode

	// Version returns the HTTP version
	Version() Version

	// Headers returns the response headers
	Headers() Header

	// Body returns the response body reader
	Body() io.Reader

	// SetStatusCode sets the HTTP status code
	SetStatusCode(StatusCode)

	// SetVersion sets the HTTP version
	SetVersion(Version)

	// SetHeader sets a header value
	SetHeader(string, string)

	// AddHeader adds a header value
	AddHeader(string, string)

	// SetBody sets the response body
	SetBody(io.Reader)

	// ContentLength returns the content length
	ContentLength() int64

	// WriteTo writes the response to a writer
	WriteTo(io.Writer) (int64, error)

	// GetHeader returns the first value of the header
	GetHeader(string) string

	// GetHeaders returns all values for the header
	GetHeaders(string) []string

	// HasHeader checks if a header exists
	HasHeader(string) bool
}

// RequestParser parses HTTP requests from raw data
type RequestParser interface {
	// Parse parses an HTTP request from a reader
	Parse(io.Reader) (Request, error)

	// ParseWithTimeout parses with a timeout
	ParseWithTimeout(io.Reader, time.Duration) (Request, error)

	// ParseBytes parses from byte slice
	ParseBytes([]byte) (Request, error)

	// Validate validates the parsed request
	Validate(Request) error
}

// ResponseBuilder builds HTTP responses
type ResponseBuilder interface {
	// Build builds an HTTP response
	Build(StatusCode, Header, io.Reader) Response

	// BuildText builds a text response
	BuildText(StatusCode, string) Response

	// BuildJSON builds a JSON response
	BuildJSON(StatusCode, interface{}) Response

	// BuildError builds an error response
	BuildError(StatusCode, string) Response

	// BuildFile builds a file response
	BuildFile(StatusCode, string) Response
}

// RequestHandler handles HTTP requests
type RequestHandler func(Request) Response

// MiddlewareFunc represents middleware function
type MiddlewareFunc func(RequestHandler) RequestHandler

// Router handles request routing
type Router interface {
	// Handle registers a handler for a method and path
	Handle(Method, string, RequestHandler)

	// HandleFunc registers a handler function
	HandleFunc(Method, string, func(Request) Response)

	// Use adds middleware
	Use(MiddlewareFunc)

	// Route finds the appropriate handler for a request
	Route(Request) (RequestHandler, map[string]string)

	// ServeRequest serves an HTTP request
	ServeRequest(Request) Response
}

// Server represents an HTTP server
type Server interface {
	// Start starts the HTTP server
	Start() error

	// Stop stops the HTTP server
	Stop() error

	// IsRunning returns true if the server is running
	IsRunning() bool

	// Addr returns the server's listening address
	Addr() net.Addr

	// SetRouter sets the request router
	SetRouter(Router)

	// SetHandler sets a single request handler
	SetHandler(RequestHandler)

	// SetMiddleware adds middleware
	SetMiddleware(...MiddlewareFunc)
}

// Client represents an HTTP client
type Client interface {
	// Get sends a GET request
	Get(string) (Response, error)

	// Post sends a POST request
	Post(string, io.Reader) (Response, error)

	// Put sends a PUT request
	Put(string, io.Reader) (Response, error)

	// Delete sends a DELETE request
	Delete(string) (Response, error)

	// Do sends a custom request
	Do(Request) (Response, error)

	// SetTimeout sets the request timeout
	SetTimeout(time.Duration)

	// SetHeader sets a default header
	SetHeader(string, string)
}

// MessageWriter writes HTTP messages to connections
type MessageWriter interface {
	// WriteRequest writes an HTTP request
	WriteRequest(io.Writer, Request) error

	// WriteResponse writes an HTTP response
	WriteResponse(io.Writer, Response) error

	// WriteHeaders writes HTTP headers
	WriteHeaders(io.Writer, Header) error

	// WriteStatusLine writes HTTP status line
	WriteStatusLine(io.Writer, Version, StatusCode) error
}

// MessageReader reads HTTP messages from connections
type MessageReader interface {
	// ReadRequest reads an HTTP request
	ReadRequest(io.Reader) (Request, error)

	// ReadResponse reads an HTTP response
	ReadResponse(io.Reader) (Response, error)

	// ReadHeaders reads HTTP headers
	ReadHeaders(io.Reader) (Header, error)

	// ReadStatusLine reads HTTP status line
	ReadStatusLine(io.Reader) (Version, StatusCode, string, error)
}

// RequestValidator validates HTTP requests
type RequestValidator interface {
	// ValidateMethod validates HTTP method
	ValidateMethod(Method) error

	// ValidatePath validates request path
	ValidatePath(string) error

	// ValidateHeaders validates headers
	ValidateHeaders(Header) error

	// ValidateVersion validates HTTP version
	ValidateVersion(Version) error

	// ValidateRequest validates complete request
	ValidateRequest(Request) error
}

// ResponseValidator validates HTTP responses
type ResponseValidator interface {
	// ValidateStatusCode validates status code
	ValidateStatusCode(StatusCode) error

	// ValidateHeaders validates headers
	ValidateHeaders(Header) error

	// ValidateVersion validates HTTP version
	ValidateVersion(Version) error

	// ValidateResponse validates complete response
	ValidateResponse(Response) error
}