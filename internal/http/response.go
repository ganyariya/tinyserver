package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ganyariya/tinyserver/internal/common"
	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

// responseImpl provides internal implementation for HTTP responses
type responseImpl struct {
	*pkghttp.httpResponse
}

// NewResponseFromRaw creates a response from raw HTTP data
func NewResponseFromRaw(rawData []byte) (pkghttp.Response, error) {
	reader := bytes.NewReader(rawData)
	return ParseResponse(reader)
}

// ParseResponse parses an HTTP response from a reader
func ParseResponse(r io.Reader) (pkghttp.Response, error) {
	scanner := bufio.NewScanner(r)
	
	// Parse status line
	if !scanner.Scan() {
		return nil, common.HTTPError("invalid response status line")
	}
	
	statusLine := scanner.Text()
	version, statusCode, err := parseStatusLine(statusLine)
	if err != nil {
		return nil, err
	}
	
	// Create response
	resp := pkghttp.NewResponse(statusCode, version).(*pkghttp.httpResponse)
	
	// Parse headers
	headers, err := parseResponseHeaders(scanner)
	if err != nil {
		return nil, err
	}
	
	// Set headers
	for name, values := range headers {
		for _, value := range values {
			resp.AddHeader(name, value)
		}
	}
	
	// Parse body if present
	contentLength := resp.ContentLength()
	if contentLength > 0 {
		body, err := parseResponseBody(scanner, contentLength)
		if err != nil {
			return nil, err
		}
		resp.SetBody(bytes.NewReader(body))
	}
	
	return resp, nil
}

// parseStatusLine parses the HTTP status line
func parseStatusLine(line string) (pkghttp.Version, pkghttp.StatusCode, error) {
	if line == "" {
		return "", 0, common.HTTPError("empty status line")
	}
	
	// Split status line into components
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		return "", 0, common.HTTPError("invalid status line format")
	}
	
	versionStr := parts[0]
	statusCodeStr := parts[1]
	// parts[2] would be the reason phrase (optional)
	
	// Validate version
	version := pkghttp.Version(versionStr)
	if !isValidVersion(version) {
		return "", 0, common.HTTPError(ErrInvalidVersion)
	}
	
	// Parse status code
	statusCodeInt, err := strconv.Atoi(statusCodeStr)
	if err != nil || statusCodeInt < 100 || statusCodeInt >= 600 {
		return "", 0, common.HTTPError(ErrInvalidStatusCode)
	}
	
	statusCode := pkghttp.StatusCode(statusCodeInt)
	
	return version, statusCode, nil
}

// parseResponseHeaders parses HTTP response headers
func parseResponseHeaders(scanner *bufio.Scanner) (pkghttp.Header, error) {
	return parseHeaders(scanner)
}

// parseResponseBody parses the response body
func parseResponseBody(scanner *bufio.Scanner, contentLength int64) ([]byte, error) {
	return parseBody(scanner, contentLength)
}

// WriteResponse writes an HTTP response to a writer
func WriteResponse(w io.Writer, resp pkghttp.Response) error {
	// Write status line
	statusLine := fmt.Sprintf("%s %d %s\r\n", 
		resp.Version(), 
		resp.StatusCode(), 
		pkghttp.StatusText(resp.StatusCode()))
	
	if _, err := w.Write([]byte(statusLine)); err != nil {
		return common.HTTPError("failed to write status line")
	}
	
	// Write headers
	for name, values := range resp.Headers() {
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
	if resp.Body() != nil {
		if _, err := io.Copy(w, resp.Body()); err != nil {
			return common.HTTPError("failed to write body")
		}
	}
	
	return nil
}

// BuildErrorResponse builds a standard error response
func BuildErrorResponse(statusCode pkghttp.StatusCode, message string) pkghttp.Response {
	if message == "" {
		message = pkghttp.StatusText(statusCode)
	}
	
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%d %s</title>
</head>
<body>
    <h1>%d %s</h1>
    <p>%s</p>
    <hr>
    <p><em>TinyServer</em></p>
</body>
</html>`, 
		statusCode, pkghttp.StatusText(statusCode),
		statusCode, pkghttp.StatusText(statusCode),
		message)
	
	return pkghttp.NewHTMLResponse(statusCode, pkghttp.Version11, html)
}

// BuildJSONErrorResponse builds a JSON error response
func BuildJSONErrorResponse(statusCode pkghttp.StatusCode, message string) pkghttp.Response {
	if message == "" {
		message = pkghttp.StatusText(statusCode)
	}
	
	json := fmt.Sprintf(`{
    "error": {
        "code": %d,
        "message": "%s"
    }
}`, statusCode, message)
	
	return pkghttp.NewJSONResponse(statusCode, pkghttp.Version11, json)
}

// BuildTextResponse builds a simple text response
func BuildTextResponse(statusCode pkghttp.StatusCode, text string) pkghttp.Response {
	return pkghttp.NewTextResponse(statusCode, pkghttp.Version11, text)
}

// BuildRedirectResponse builds a redirect response
func BuildRedirectResponse(statusCode pkghttp.StatusCode, location string) pkghttp.Response {
	resp := pkghttp.NewResponse(statusCode, pkghttp.Version11)
	resp.SetHeader(pkghttp.HeaderLocation, location)
	
	// Add redirect body
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%d %s</title>
</head>
<body>
    <h1>%d %s</h1>
    <p>The document has moved <a href="%s">here</a>.</p>
    <hr>
    <p><em>TinyServer</em></p>
</body>
</html>`, 
		statusCode, pkghttp.StatusText(statusCode),
		statusCode, pkghttp.StatusText(statusCode),
		location)
	
	resp.SetBody(strings.NewReader(html))
	resp.SetHeader(pkghttp.HeaderContentType, pkghttp.MimeTypeTextHTML)
	resp.SetHeader(pkghttp.HeaderContentLength, strconv.Itoa(len(html)))
	
	return resp
}

// FormatResponse formats a response for debugging/logging
func FormatResponse(resp pkghttp.Response) string {
	var buf bytes.Buffer
	
	// Status line
	fmt.Fprintf(&buf, "%s %d %s\n", 
		resp.Version(), 
		resp.StatusCode(), 
		pkghttp.StatusText(resp.StatusCode()))
	
	// Headers
	for name, values := range resp.Headers() {
		for _, value := range values {
			fmt.Fprintf(&buf, "%s: %s\n", name, value)
		}
	}
	
	// Content length
	if contentLength := resp.ContentLength(); contentLength > 0 {
		fmt.Fprintf(&buf, "Content-Length: %d\n", contentLength)
	}
	
	// Status information
	fmt.Fprintf(&buf, "Success: %t\n", resp.IsSuccess())
	fmt.Fprintf(&buf, "Error: %t\n", resp.IsError())
	fmt.Fprintf(&buf, "Redirection: %t\n", resp.IsRedirection())
	
	return buf.String()
}

// SetCommonHeaders sets common response headers
func SetCommonHeaders(resp pkghttp.Response) {
	// Set server header
	resp.SetHeader(pkghttp.HeaderServer, "TinyServer/1.0")
	
	// Set date header
	resp.SetHeader(pkghttp.HeaderDate, common.FormatHTTPDate())
	
	// Set connection header (default to close for simplicity)
	resp.SetHeader(pkghttp.HeaderConnection, "close")
}

// ValidateResponse validates a response
func ValidateResponse(resp pkghttp.Response) error {
	// Validate status code
	if resp.StatusCode() < 100 || resp.StatusCode() >= 600 {
		return common.HTTPError(ErrInvalidStatusCode)
	}
	
	// Validate version
	if !isValidVersion(resp.Version()) {
		return common.HTTPError(ErrInvalidVersion)
	}
	
	// Validate headers
	for name := range resp.Headers() {
		if !isValidHeaderName(name) {
			return common.HTTPError(ErrInvalidHeader)
		}
	}
	
	// Validate content length consistency
	contentLength := resp.ContentLength()
	if contentLength < 0 {
		return common.HTTPError(ErrInvalidContentLength)
	}
	
	return nil
}