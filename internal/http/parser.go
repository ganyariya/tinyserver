package http

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"time"

	"github.com/ganyariya/tinyserver/internal/common"
	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

// httpParser implements HTTP parsing functionality
type httpParser struct {
	logger *common.Logger
}

// NewParser creates a new HTTP parser
func NewParser() pkghttp.RequestParser {
	return &httpParser{
		logger: common.NewDefaultLogger(),
	}
}

// Parse parses an HTTP request from a reader
func (p *httpParser) Parse(r io.Reader) (pkghttp.Request, error) {
	return ParseRequest(r, nil)
}

// ParseWithTimeout parses with a timeout
func (p *httpParser) ParseWithTimeout(r io.Reader, timeout time.Duration) (pkghttp.Request, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a channel to receive the result
	resultChan := make(chan parseResult, 1)

	// Parse in a goroutine
	go func() {
		req, err := p.Parse(r)
		resultChan <- parseResult{req: req, err: err}
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result.req, result.err
	case <-ctx.Done():
		return nil, common.TimeoutError(ErrParseTimeout)
	}
}

// ParseBytes parses from byte slice
func (p *httpParser) ParseBytes(data []byte) (pkghttp.Request, error) {
	reader := bytes.NewReader(data)
	return p.Parse(reader)
}

// Validate validates the parsed request
func (p *httpParser) Validate(req pkghttp.Request) error {
	if req == nil {
		return common.HTTPError("request is nil")
	}

	// Validate method
	if req.Method() == "" {
		return common.HTTPError(ErrInvalidMethod)
	}

	if !isValidMethod(req.Method()) {
		return common.HTTPError(ErrInvalidMethod)
	}

	// Validate path
	if req.Path() == "" {
		return common.HTTPError(ErrInvalidPath)
	}

	if !isValidPath(req.Path()) {
		return common.HTTPError(ErrInvalidPath)
	}

	// Validate version
	if !isValidVersion(req.Version()) {
		return common.HTTPError(ErrInvalidVersion)
	}

	// Validate headers
	for name := range req.Headers() {
		if !isValidHeaderName(name) {
			return common.HTTPError(ErrInvalidHeader)
		}
	}

	// Validate content length consistency
	contentLength := req.ContentLength()
	if contentLength < 0 {
		return common.HTTPError(ErrInvalidContentLength)
	}

	return nil
}

// parseResult holds the result of parsing operation
type parseResult struct {
	req pkghttp.Request
	err error
}

// httpResponseParser implements HTTP response parsing functionality
type httpResponseParser struct {
	logger *common.Logger
}

// NewResponseParser creates a new HTTP response parser
func NewResponseParser() *httpResponseParser {
	return &httpResponseParser{
		logger: common.NewDefaultLogger(),
	}
}

// ParseResponse parses an HTTP response with timeout
func (p *httpResponseParser) ParseResponse(r io.Reader) (pkghttp.Response, error) {
	return ParseResponse(r)
}

// ParseResponseWithTimeout parses a response with timeout
func (p *httpResponseParser) ParseResponseWithTimeout(r io.Reader, timeout time.Duration) (pkghttp.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a channel to receive the result
	resultChan := make(chan responseParseResult, 1)

	// Parse in a goroutine
	go func() {
		resp, err := p.ParseResponse(r)
		resultChan <- responseParseResult{resp: resp, err: err}
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result.resp, result.err
	case <-ctx.Done():
		return nil, common.TimeoutError(ErrParseTimeout)
	}
}

// ParseResponseBytes parses response from byte slice
func (p *httpResponseParser) ParseResponseBytes(data []byte) (pkghttp.Response, error) {
	reader := bytes.NewReader(data)
	return p.ParseResponse(reader)
}

// responseParseResult holds the result of response parsing operation
type responseParseResult struct {
	resp pkghttp.Response
	err  error
}

// messageParser provides unified parsing for HTTP messages
type messageParser struct {
	logger        *common.Logger
	maxHeaderSize int
	maxBodySize   int64
}

// NewMessageParser creates a new message parser
func NewMessageParser() *messageParser {
	return &messageParser{
		logger:        common.NewDefaultLogger(),
		maxHeaderSize: pkghttp.MaxHeaderSize,
		maxBodySize:   pkghttp.MaxRequestBodySize,
	}
}

// SetMaxHeaderSize sets the maximum header size
func (p *messageParser) SetMaxHeaderSize(size int) {
	p.maxHeaderSize = size
}

// SetMaxBodySize sets the maximum body size
func (p *messageParser) SetMaxBodySize(size int64) {
	p.maxBodySize = size
}

// ParseHTTPMessage parses a generic HTTP message
func (p *messageParser) ParseHTTPMessage(r io.Reader) ([]string, pkghttp.Header, io.Reader, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	var totalSize int

	// Read until we find the first line (status/request line)
	if !scanner.Scan() {
		return nil, nil, nil, common.HTTPError(ErrUnexpectedEOF)
	}

	firstLine := scanner.Text()
	lines = append(lines, firstLine)
	totalSize += len(firstLine)

	if totalSize > p.maxHeaderSize {
		return nil, nil, nil, common.HTTPError(ErrHeaderTooLarge)
	}

	// Parse headers
	headers, err := parseHeaders(scanner)
	if err != nil {
		return nil, nil, nil, err
	}

	// Calculate remaining data for body
	var bodyReader io.Reader
	if scanner.Scan() {
		// If there's more data, create a reader for the body
		firstBodyLine := scanner.Bytes()
		bodyReader = io.MultiReader(
			bytes.NewReader(firstBodyLine),
			bytes.NewReader([]byte("\n")),
			r, // Original reader for remaining data
		)
	}

	return lines, headers, bodyReader, nil
}

// ChunkedReader handles chunked transfer encoding
type ChunkedReader struct {
	r      *bufio.Reader
	n      int64 // bytes remaining in current chunk
	err    error
	logger *common.Logger
}

// NewChunkedReader creates a new chunked reader
func NewChunkedReader(r io.Reader) *ChunkedReader {
	return &ChunkedReader{
		r:      bufio.NewReader(r),
		logger: common.NewDefaultLogger(),
	}
}

// Read implements io.Reader for chunked data
func (cr *ChunkedReader) Read(p []byte) (int, error) {
	if cr.err != nil {
		return 0, cr.err
	}

	if cr.n == 0 {
		// Read next chunk size
		line, _, err := cr.r.ReadLine()
		if err != nil {
			cr.err = err
			return 0, err
		}

		// Parse chunk size (hexadecimal)
		chunkSize, err := parseChunkSize(string(line))
		if err != nil {
			cr.err = common.HTTPError(ErrChunkedEncodingInvalid)
			return 0, cr.err
		}

		if chunkSize == 0 {
			// End of chunks, read trailing headers if any
			cr.readTrailers()
			cr.err = io.EOF
			return 0, io.EOF
		}

		if chunkSize > MaxChunkSize {
			cr.err = common.HTTPError(ErrChunkedEncodingInvalid)
			return 0, cr.err
		}

		cr.n = int64(chunkSize)
	}

	// Read data from current chunk
	if int64(len(p)) > cr.n {
		p = p[:cr.n]
	}

	n, err := cr.r.Read(p)
	cr.n -= int64(n)

	if cr.n == 0 && err == nil {
		// End of chunk, read trailing CRLF
		cr.r.ReadLine()
	}

	if err != nil {
		cr.err = err
	}

	return n, err
}

// readTrailers reads any trailing headers after the last chunk
func (cr *ChunkedReader) readTrailers() {
	// Read trailing headers (usually empty)
	for {
		line, _, err := cr.r.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}
		// Log any trailing headers for debugging
		cr.logger.Debug("Trailing header: %s", string(line))
	}
}

// parseChunkSize parses hexadecimal chunk size
func parseChunkSize(line string) (int, error) {
	// Remove any chunk extensions (after semicolon)
	if idx := bytes.IndexByte([]byte(line), ';'); idx >= 0 {
		line = line[:idx]
	}

	// Parse hexadecimal
	var size int
	for _, b := range []byte(line) {
		if b >= '0' && b <= '9' {
			size = size*16 + int(b-'0')
		} else if b >= 'a' && b <= 'f' {
			size = size*16 + int(b-'a'+10)
		} else if b >= 'A' && b <= 'F' {
			size = size*16 + int(b-'A'+10)
		} else {
			return 0, common.HTTPError(ErrChunkedEncodingInvalid)
		}
	}

	return size, nil
}

// ContentLengthReader handles content-length based reading
type ContentLengthReader struct {
	r         io.Reader
	remaining int64
	logger    *common.Logger
}

// NewContentLengthReader creates a new content-length reader
func NewContentLengthReader(r io.Reader, contentLength int64) *ContentLengthReader {
	return &ContentLengthReader{
		r:         r,
		remaining: contentLength,
		logger:    common.NewDefaultLogger(),
	}
}

// Read implements io.Reader for content-length based reading
func (clr *ContentLengthReader) Read(p []byte) (int, error) {
	if clr.remaining <= 0 {
		return 0, io.EOF
	}

	if int64(len(p)) > clr.remaining {
		p = p[:clr.remaining]
	}

	n, err := clr.r.Read(p)
	clr.remaining -= int64(n)

	return n, err
}

// Remaining returns the number of bytes remaining
func (clr *ContentLengthReader) Remaining() int64 {
	return clr.remaining
}

// HTTPMessageBuilder helps build HTTP messages
type HTTPMessageBuilder struct {
	logger *common.Logger
}

// NewHTTPMessageBuilder creates a new message builder
func NewHTTPMessageBuilder() *HTTPMessageBuilder {
	return &HTTPMessageBuilder{
		logger: common.NewDefaultLogger(),
	}
}

// BuildRequest builds an HTTP request message
func (b *HTTPMessageBuilder) BuildRequest(req pkghttp.Request) ([]byte, error) {
	var buf bytes.Buffer

	if err := WriteRequest(&buf, req); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// BuildResponse builds an HTTP response message
func (b *HTTPMessageBuilder) BuildResponse(resp pkghttp.Response) ([]byte, error) {
	var buf bytes.Buffer

	if err := WriteResponse(&buf, resp); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
