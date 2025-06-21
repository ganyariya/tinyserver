package http

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// httpResponse implements the Response interface
type httpResponse struct {
	statusCode StatusCode
	version    Version
	headers    Header
	body       io.Reader
}

// NewResponse creates a new HTTP response
func NewResponse(statusCode StatusCode, version Version) Response {
	return &httpResponse{
		statusCode: statusCode,
		version:    version,
		headers:    make(Header),
	}
}

// NewResponseWithBody creates a new HTTP response with body
func NewResponseWithBody(statusCode StatusCode, version Version, body io.Reader) Response {
	resp := &httpResponse{
		statusCode: statusCode,
		version:    version,
		headers:    make(Header),
		body:       body,
	}
	return resp
}

// NewTextResponse creates a text response
func NewTextResponse(statusCode StatusCode, version Version, text string) Response {
	resp := &httpResponse{
		statusCode: statusCode,
		version:    version,
		headers:    make(Header),
		body:       strings.NewReader(text),
	}
	
	// Set content type and length
	resp.SetHeader(HeaderContentType, MimeTypeTextPlain)
	resp.SetHeader(HeaderContentLength, strconv.Itoa(len(text)))
	
	return resp
}

// NewHTMLResponse creates an HTML response
func NewHTMLResponse(statusCode StatusCode, version Version, html string) Response {
	resp := &httpResponse{
		statusCode: statusCode,
		version:    version,
		headers:    make(Header),
		body:       strings.NewReader(html),
	}
	
	// Set content type and length
	resp.SetHeader(HeaderContentType, MimeTypeTextHTML)
	resp.SetHeader(HeaderContentLength, strconv.Itoa(len(html)))
	
	return resp
}

// NewJSONResponse creates a JSON response
func NewJSONResponse(statusCode StatusCode, version Version, json string) Response {
	resp := &httpResponse{
		statusCode: statusCode,
		version:    version,
		headers:    make(Header),
		body:       strings.NewReader(json),
	}
	
	// Set content type and length
	resp.SetHeader(HeaderContentType, MimeTypeJSON)
	resp.SetHeader(HeaderContentLength, strconv.Itoa(len(json)))
	
	return resp
}

// StatusCode returns the HTTP status code
func (r *httpResponse) StatusCode() StatusCode {
	return r.statusCode
}

// Version returns the HTTP version
func (r *httpResponse) Version() Version {
	return r.version
}

// Headers returns the response headers
func (r *httpResponse) Headers() Header {
	if r.headers == nil {
		r.headers = make(Header)
	}
	return r.headers
}

// Body returns the response body reader
func (r *httpResponse) Body() io.Reader {
	return r.body
}

// SetStatusCode sets the HTTP status code
func (r *httpResponse) SetStatusCode(statusCode StatusCode) {
	r.statusCode = statusCode
}

// SetVersion sets the HTTP version
func (r *httpResponse) SetVersion(version Version) {
	r.version = version
}

// SetHeader sets a header value
func (r *httpResponse) SetHeader(name, value string) {
	if r.headers == nil {
		r.headers = make(Header)
	}
	r.headers[name] = []string{value}
}

// AddHeader adds a header value
func (r *httpResponse) AddHeader(name, value string) {
	if r.headers == nil {
		r.headers = make(Header)
	}
	r.headers[name] = append(r.headers[name], value)
}

// SetBody sets the response body
func (r *httpResponse) SetBody(body io.Reader) {
	r.body = body
}

// ContentLength returns the content length
func (r *httpResponse) ContentLength() int64 {
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

// WriteTo writes the response to a writer
func (r *httpResponse) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64
	
	// Write status line
	statusLine := fmt.Sprintf("%s %d %s%s", 
		r.version, 
		r.statusCode, 
		StatusText(r.statusCode),
		HTTPSeparator)
	
	n, err := w.Write([]byte(statusLine))
	totalWritten += int64(n)
	if err != nil {
		return totalWritten, err
	}
	
	// Write headers
	if r.headers != nil {
		for name, values := range r.headers {
			for _, value := range values {
				headerLine := fmt.Sprintf("%s%s%s%s", 
					name, 
					HTTPHeaderSeparator, 
					value,
					HTTPSeparator)
				
				n, err := w.Write([]byte(headerLine))
				totalWritten += int64(n)
				if err != nil {
					return totalWritten, err
				}
			}
		}
	}
	
	// Write header-body separator
	n, err = w.Write([]byte(HTTPSeparator))
	totalWritten += int64(n)
	if err != nil {
		return totalWritten, err
	}
	
	// Write body if present
	if r.body != nil {
		n, err := io.Copy(w, r.body)
		totalWritten += n
		if err != nil {
			return totalWritten, err
		}
	}
	
	return totalWritten, nil
}

// GetHeader returns the first value of the header
func (r *httpResponse) GetHeader(name string) string {
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
func (r *httpResponse) GetHeaders(name string) []string {
	if r.headers == nil {
		return nil
	}
	
	return r.headers[name]
}

// HasHeader checks if a header exists
func (r *httpResponse) HasHeader(name string) bool {
	if r.headers == nil {
		return false
	}
	
	_, exists := r.headers[name]
	return exists
}

// SetContentType sets the Content-Type header
func (r *httpResponse) SetContentType(contentType string) {
	r.SetHeader(HeaderContentType, contentType)
}

// SetContentLength sets the Content-Length header
func (r *httpResponse) SetContentLength(length int64) {
	r.SetHeader(HeaderContentLength, strconv.FormatInt(length, 10))
}

// String returns a string representation of the response
func (r *httpResponse) String() string {
	var buf bytes.Buffer
	r.WriteTo(&buf)
	return buf.String()
}

// Clone creates a copy of the response
func (r *httpResponse) Clone() Response {
	clone := &httpResponse{
		statusCode: r.statusCode,
		version:    r.version,
		headers:    make(Header),
		body:       r.body,
	}
	
	// Deep copy headers
	for name, values := range r.headers {
		clone.headers[name] = make([]string, len(values))
		copy(clone.headers[name], values)
	}
	
	return clone
}

// IsSuccess returns true if the status code indicates success
func (r *httpResponse) IsSuccess() bool {
	return IsSuccess(r.statusCode)
}

// IsError returns true if the status code indicates an error
func (r *httpResponse) IsError() bool {
	return IsError(r.statusCode)
}

// IsRedirection returns true if the status code indicates redirection
func (r *httpResponse) IsRedirection() bool {
	return IsRedirection(r.statusCode)
}