package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/ganyariya/tinyserver/internal/common"
	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

// requestImpl provides internal implementation for HTTP requests
type requestImpl struct {
	*pkghttp.HTTPRequest
}

// NewRequestFromRaw creates a request from raw HTTP data
func NewRequestFromRaw(rawData []byte, remoteAddr net.Addr) (pkghttp.Request, error) {
	reader := bytes.NewReader(rawData)
	return ParseRequest(reader, remoteAddr)
}

// ParseRequest parses an HTTP request from a reader
func ParseRequest(r io.Reader, remoteAddr net.Addr) (pkghttp.Request, error) {
	scanner := bufio.NewScanner(r)
	
	// Parse request line
	if !scanner.Scan() {
		return nil, common.HTTPError(ErrInvalidRequestLine)
	}
	
	requestLine := scanner.Text()
	method, path, version, err := parseRequestLine(requestLine)
	if err != nil {
		return nil, err
	}
	
	// Create request
	req := pkghttp.NewRequest(method, path, version).(*pkghttp.HTTPRequest)
	req.SetRemoteAddr(remoteAddr)
	
	// Parse headers
	headers, err := parseHeaders(scanner)
	if err != nil {
		return nil, err
	}
	
	// Set headers
	for name, values := range headers {
		for _, value := range values {
			req.AddHeader(name, value)
		}
	}
	
	// Parse body if present
	contentLength := req.ContentLength()
	if contentLength > 0 {
		body, err := parseBody(scanner, contentLength)
		if err != nil {
			return nil, err
		}
		req.SetBody(bytes.NewReader(body))
	}
	
	return req, nil
}

// parseRequestLine parses the HTTP request line
func parseRequestLine(line string) (pkghttp.Method, string, pkghttp.Version, error) {
	if line == "" {
		return "", "", "", common.HTTPError(ErrInvalidRequestLine)
	}
	
	if len(line) > MaxRequestLineLength {
		return "", "", "", common.HTTPError(ErrRequestTooLarge)
	}
	
	// Split request line into components
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return "", "", "", common.HTTPError(ErrInvalidRequestLine)
	}
	
	methodStr := parts[0]
	path := parts[1]
	versionStr := parts[2]
	
	// Validate method
	method := pkghttp.Method(methodStr)
	if !isValidMethod(method) {
		return "", "", "", common.HTTPError(ErrInvalidMethod)
	}
	
	// Validate path
	if !isValidPath(path) {
		return "", "", "", common.HTTPError(ErrInvalidPath)
	}
	
	// Validate version
	version := pkghttp.Version(versionStr)
	if !isValidVersion(version) {
		return "", "", "", common.HTTPError(ErrInvalidVersion)
	}
	
	return method, path, version, nil
}

// parseHeaders parses HTTP headers
func parseHeaders(scanner *bufio.Scanner) (pkghttp.Header, error) {
	headers := make(pkghttp.Header)
	headerCount := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Empty line indicates end of headers
		if line == "" {
			break
		}
		
		// Check header count limit
		headerCount++
		if headerCount > MaxHeaderLines {
			return nil, common.HTTPError(ErrHeaderTooLarge)
		}
		
		// Check line length
		if len(line) > MaxHeaderLineLength {
			return nil, common.HTTPError(ErrHeaderTooLarge)
		}
		
		// Parse header
		name, value, err := parseHeader(line)
		if err != nil {
			return nil, err
		}
		
		headers[name] = append(headers[name], value)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, common.HTTPError(ErrUnexpectedEOF)
	}
	
	return headers, nil
}

// parseHeader parses a single header line
func parseHeader(line string) (string, string, error) {
	// Find colon separator
	colonIndex := strings.Index(line, ":")
	if colonIndex == -1 {
		return "", "", common.HTTPError(ErrInvalidHeader)
	}
	
	name := strings.TrimSpace(line[:colonIndex])
	value := strings.TrimSpace(line[colonIndex+1:])
	
	// Validate header name
	if !isValidHeaderName(name) {
		return "", "", common.HTTPError(ErrInvalidHeader)
	}
	
	return name, value, nil
}

// parseBody parses the request body
func parseBody(scanner *bufio.Scanner, contentLength int64) ([]byte, error) {
	if contentLength <= 0 {
		return nil, nil
	}
	
	if contentLength > pkghttp.MaxRequestBodySize {
		return nil, common.HTTPError(ErrRequestTooLarge)
	}
	
	body := make([]byte, 0, contentLength)
	remainingBytes := contentLength
	
	for scanner.Scan() && remainingBytes > 0 {
		line := scanner.Bytes()
		
		// Add line to body
		if int64(len(line)) <= remainingBytes {
			body = append(body, line...)
			remainingBytes -= int64(len(line))
		} else {
			// Partial line
			body = append(body, line[:remainingBytes]...)
			remainingBytes = 0
		}
		
		// Add newline if not last line
		if remainingBytes > 0 {
			body = append(body, '\n')
			remainingBytes--
		}
	}
	
	if remainingBytes > 0 {
		return nil, common.HTTPError(ErrUnexpectedEOF)
	}
	
	return body, nil
}

// Validation functions

// isValidMethod checks if the method is valid
func isValidMethod(method pkghttp.Method) bool {
	switch method {
	case pkghttp.MethodGet, pkghttp.MethodPost, pkghttp.MethodPut, 
		 pkghttp.MethodDelete, pkghttp.MethodHead, pkghttp.MethodOptions, 
		 pkghttp.MethodPatch:
		return true
	default:
		return false
	}
}

// isValidPath checks if the path is valid
func isValidPath(path string) bool {
	if path == "" {
		return false
	}
	
	// Path must start with /
	if !strings.HasPrefix(path, "/") {
		return false
	}
	
	// Basic validation - no control characters
	for _, r := range path {
		if r < 32 || r == 127 {
			return false
		}
	}
	
	return true
}

// isValidVersion checks if the HTTP version is valid
func isValidVersion(version pkghttp.Version) bool {
	switch version {
	case pkghttp.Version10, pkghttp.Version11:
		return true
	default:
		return false
	}
}

// isValidHeaderName checks if the header name is valid
func isValidHeaderName(name string) bool {
	if name == "" {
		return false
	}
	
	// Header names can contain letters, digits, and hyphens
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || 
			 (r >= 'A' && r <= 'Z') || 
			 (r >= '0' && r <= '9') || 
			 r == '-') {
			return false
		}
	}
	
	return true
}

// WriteRequest writes an HTTP request to a writer
func WriteRequest(w io.Writer, req pkghttp.Request) error {
	// Write request line
	requestLine := fmt.Sprintf("%s %s %s\r\n", 
		req.Method(), 
		req.Path(), 
		req.Version())
	
	if _, err := w.Write([]byte(requestLine)); err != nil {
		return common.HTTPError("failed to write request line")
	}
	
	// Write headers
	for name, values := range req.Headers() {
		for _, value := range values {
			headerLine := fmt.Sprintf("%s: %s\r\n", name, value)
			if _, err := w.Write([]byte(headerLine)); err != nil {
				return common.HTTPError("failed to write header")
			}
		}
	}
	
	// Write header-body separator
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return common.HTTPError("failed to write header separator")
	}
	
	// Write body if present
	if req.Body() != nil {
		if _, err := io.Copy(w, req.Body()); err != nil {
			return common.HTTPError("failed to write body")
		}
	}
	
	return nil
}

// FormatRequest formats a request for debugging/logging
func FormatRequest(req pkghttp.Request) string {
	var buf bytes.Buffer
	
	// Request line
	fmt.Fprintf(&buf, "%s %s %s\n", req.Method(), req.Path(), req.Version())
	
	// Headers
	for name, values := range req.Headers() {
		for _, value := range values {
			fmt.Fprintf(&buf, "%s: %s\n", name, value)
		}
	}
	
	// Remote address if available
	if req.RemoteAddr() != nil {
		fmt.Fprintf(&buf, "Remote-Addr: %s\n", req.RemoteAddr().String())
	}
	
	// Content length
	if contentLength := req.ContentLength(); contentLength > 0 {
		fmt.Fprintf(&buf, "Content-Length: %d\n", contentLength)
	}
	
	return buf.String()
}