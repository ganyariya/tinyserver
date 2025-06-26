package http

import (
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// HTTPRequest implements the Request interface
type HTTPRequest struct {
	method      Method
	path        string
	version     Version
	headers     Header
	body        io.Reader
	queryParams map[string]string
	remoteAddr  net.Addr
}

// NewRequest creates a new HTTP request
func NewRequest(method Method, path string, version Version) Request {
	return &HTTPRequest{
		method:      method,
		path:        path,
		version:     version,
		headers:     make(Header),
		queryParams: make(map[string]string),
	}
}

// NewRequestWithBody creates a new HTTP request with body
func NewRequestWithBody(method Method, path string, version Version, body io.Reader) Request {
	req := &HTTPRequest{
		method:      method,
		path:        path,
		version:     version,
		headers:     make(Header),
		queryParams: make(map[string]string),
		body:        body,
	}
	return req
}

// Method returns the HTTP method
func (r *HTTPRequest) Method() Method {
	return r.method
}

// Path returns the request path
func (r *HTTPRequest) Path() string {
	return r.path
}

// Version returns the HTTP version
func (r *HTTPRequest) Version() Version {
	return r.version
}

// Headers returns the request headers
func (r *HTTPRequest) Headers() Header {
	if r.headers == nil {
		r.headers = make(Header)
	}
	return r.headers
}

// Body returns the request body reader
func (r *HTTPRequest) Body() io.Reader {
	return r.body
}

// QueryParams returns query parameters
func (r *HTTPRequest) QueryParams() map[string]string {
	if r.queryParams == nil {
		r.queryParams = make(map[string]string)
		r.parseQueryParams()
	}
	return r.queryParams
}

// SetMethod sets the HTTP method
func (r *HTTPRequest) SetMethod(method Method) {
	r.method = method
}

// SetPath sets the request path
func (r *HTTPRequest) SetPath(path string) {
	r.path = path
	// Re-parse query parameters when path changes
	r.queryParams = make(map[string]string)
	r.parseQueryParams()
}

// SetVersion sets the HTTP version
func (r *HTTPRequest) SetVersion(version Version) {
	r.version = version
}

// SetHeader sets a header value
func (r *HTTPRequest) SetHeader(name, value string) {
	if r.headers == nil {
		r.headers = make(Header)
	}
	r.headers[name] = []string{value}
}

// AddHeader adds a header value
func (r *HTTPRequest) AddHeader(name, value string) {
	if r.headers == nil {
		r.headers = make(Header)
	}
	r.headers[name] = append(r.headers[name], value)
}

// SetBody sets the request body
func (r *HTTPRequest) SetBody(body io.Reader) {
	r.body = body
}

// ContentLength returns the content length
func (r *HTTPRequest) ContentLength() int64 {
	if r.headers == nil {
		return 0
	}

	contentLengths, exists := r.headers[HeaderContentLength]
	if !exists || len(contentLengths) == 0 {
		return 0
	}

	length, err := strconv.ParseInt(contentLengths[0], 10, 64)
	if err != nil {
		return 0
	}

	return length
}

// RemoteAddr returns the remote address
func (r *HTTPRequest) RemoteAddr() net.Addr {
	return r.remoteAddr
}

// SetRemoteAddr sets the remote address (internal method)
func (r *HTTPRequest) SetRemoteAddr(addr net.Addr) {
	r.remoteAddr = addr
}

// parseQueryParams parses query parameters from the path
func (r *HTTPRequest) parseQueryParams() {
	if r.queryParams == nil {
		r.queryParams = make(map[string]string)
	}

	if r.path == "" {
		return
	}

	// Find query string separator
	queryIndex := strings.Index(r.path, "?")
	if queryIndex == -1 {
		return
	}

	queryString := r.path[queryIndex+1:]
	if queryString == "" {
		return
	}

	// Parse query string
	params, err := url.ParseQuery(queryString)
	if err != nil {
		return
	}

	// Convert url.Values to map[string]string (take first value for each key)
	for key, values := range params {
		if len(values) > 0 {
			r.queryParams[key] = values[0]
		}
	}
}

// GetHeader returns the first value of the header
func (r *HTTPRequest) GetHeader(name string) string {
	if r.headers == nil {
		return ""
	}

	values, exists := r.headers[name]
	if !exists || len(values) == 0 {
		return ""
	}

	return values[0]
}

// GetHeaders returns all values for the header
func (r *HTTPRequest) GetHeaders(name string) []string {
	if r.headers == nil {
		return nil
	}

	return r.headers[name]
}

// HasHeader checks if a header exists
func (r *HTTPRequest) HasHeader(name string) bool {
	if r.headers == nil {
		return false
	}

	_, exists := r.headers[name]
	return exists
}

// PathWithoutQuery returns the path without query string
func (r *HTTPRequest) PathWithoutQuery() string {
	if r.path == "" {
		return ""
	}

	queryIndex := strings.Index(r.path, "?")
	if queryIndex == -1 {
		return r.path
	}

	return r.path[:queryIndex]
}

// Clone creates a copy of the request
func (r *HTTPRequest) Clone() Request {
	clone := &HTTPRequest{
		method:      r.method,
		path:        r.path,
		version:     r.version,
		headers:     make(Header),
		queryParams: make(map[string]string),
		body:        r.body,
		remoteAddr:  r.remoteAddr,
	}

	// Deep copy headers
	for name, values := range r.headers {
		clone.headers[name] = make([]string, len(values))
		copy(clone.headers[name], values)
	}

	// Deep copy query params
	for key, value := range r.queryParams {
		clone.queryParams[key] = value
	}

	return clone
}
