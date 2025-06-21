package http

import (
	"bytes"
	"strings"
	"testing"
	"time"

	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name     string
		rawData  string
		wantErr  bool
		expected struct {
			method  pkghttp.Method
			path    string
			version pkghttp.Version
		}
	}{
		{
			name: "simple GET request",
			rawData: "GET /hello HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"User-Agent: TinyClient/1.0\r\n" +
				"\r\n",
			wantErr: false,
			expected: struct {
				method  pkghttp.Method
				path    string
				version pkghttp.Version
			}{
				method:  pkghttp.MethodGet,
				path:    "/hello",
				version: pkghttp.Version11,
			},
		},
		{
			name: "POST request with body",
			rawData: "POST /api/data HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"{\"test\": true}",
			wantErr: false,
			expected: struct {
				method  pkghttp.Method
				path    string
				version pkghttp.Version
			}{
				method:  pkghttp.MethodPost,
				path:    "/api/data",
				version: pkghttp.Version11,
			},
		},
		{
			name: "invalid method",
			rawData: "INVALID /hello HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"\r\n",
			wantErr: true,
		},
		{
			name: "invalid path",
			rawData: "GET hello HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"\r\n",
			wantErr: true,
		},
		{
			name: "invalid version",
			rawData: "GET /hello HTTP/2.0\r\n" +
				"Host: example.com\r\n" +
				"\r\n",
			wantErr: true,
		},
		{
			name: "malformed request line",
			rawData: "GET /hello\r\n" +
				"Host: example.com\r\n" +
				"\r\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.rawData)
			req, err := ParseRequest(reader, nil)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if req.Method() != tt.expected.method {
				t.Errorf("Expected method %s, got %s", tt.expected.method, req.Method())
			}

			if req.Path() != tt.expected.path {
				t.Errorf("Expected path %s, got %s", tt.expected.path, req.Path())
			}

			if req.Version() != tt.expected.version {
				t.Errorf("Expected version %s, got %s", tt.expected.version, req.Version())
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name     string
		rawData  string
		wantErr  bool
		expected struct {
			statusCode pkghttp.StatusCode
			version    pkghttp.Version
		}
	}{
		{
			name: "successful response",
			rawData: "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: 12\r\n" +
				"\r\n" +
				"Hello, World!",
			wantErr: false,
			expected: struct {
				statusCode pkghttp.StatusCode
				version    pkghttp.Version
			}{
				statusCode: pkghttp.StatusOK,
				version:    pkghttp.Version11,
			},
		},
		{
			name: "error response",
			rawData: "HTTP/1.1 404 Not Found\r\n" +
				"Content-Type: text/html\r\n" +
				"Content-Length: 9\r\n" +
				"\r\n" +
				"Not Found",
			wantErr: false,
			expected: struct {
				statusCode pkghttp.StatusCode
				version    pkghttp.Version
			}{
				statusCode: pkghttp.StatusNotFound,
				version:    pkghttp.Version11,
			},
		},
		{
			name: "invalid status code",
			rawData: "HTTP/1.1 999 Unknown\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n",
			wantErr: true,
		},
		{
			name: "malformed status line",
			rawData: "HTTP/1.1\r\n" +
				"Content-Type: text/plain\r\n" +
				"\r\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.rawData)
			resp, err := ParseResponse(reader)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if resp.StatusCode() != tt.expected.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.expected.statusCode, resp.StatusCode())
			}

			if resp.Version() != tt.expected.version {
				t.Errorf("Expected version %s, got %s", tt.expected.version, resp.Version())
			}
		})
	}
}

func TestHttpParser(t *testing.T) {
	parser := NewParser()

	t.Run("parse simple request", func(t *testing.T) {
		rawData := "GET /test HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"\r\n"

		reader := strings.NewReader(rawData)
		req, err := parser.Parse(reader)

		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if req.Method() != pkghttp.MethodGet {
			t.Errorf("Expected GET, got %s", req.Method())
		}

		if req.Path() != "/test" {
			t.Errorf("Expected /test, got %s", req.Path())
		}
	})

	t.Run("parse with timeout", func(t *testing.T) {
		rawData := "GET /test HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"\r\n"

		reader := strings.NewReader(rawData)
		req, err := parser.ParseWithTimeout(reader, time.Second)

		if err != nil {
			t.Fatalf("ParseWithTimeout failed: %v", err)
		}

		if req.Method() != pkghttp.MethodGet {
			t.Errorf("Expected GET, got %s", req.Method())
		}
	})

	t.Run("parse bytes", func(t *testing.T) {
		rawData := []byte("GET /test HTTP/1.1\r\n" +
			"Host: example.com\r\n" +
			"\r\n")

		req, err := parser.ParseBytes(rawData)

		if err != nil {
			t.Fatalf("ParseBytes failed: %v", err)
		}

		if req.Method() != pkghttp.MethodGet {
			t.Errorf("Expected GET, got %s", req.Method())
		}
	})

	t.Run("validate request", func(t *testing.T) {
		req := pkghttp.NewRequest(pkghttp.MethodGet, "/test", pkghttp.Version11)
		err := parser.Validate(req)

		if err != nil {
			t.Errorf("Validate failed: %v", err)
		}
	})

	t.Run("validate invalid request", func(t *testing.T) {
		req := pkghttp.NewRequest("", "", "")
		err := parser.Validate(req)

		if err == nil {
			t.Error("Expected validation error")
		}
	})
}

func TestChunkedReader(t *testing.T) {
	t.Run("simple chunked data", func(t *testing.T) {
		chunkedData := "5\r\nHello\r\n" +
			"6\r\n World\r\n" +
			"0\r\n" +
			"\r\n"

		reader := NewChunkedReader(strings.NewReader(chunkedData))
		result := make([]byte, 20)
		n, err := reader.Read(result)

		if err != nil && err.Error() != "EOF" {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "Hello"
		if string(result[:n]) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result[:n]))
		}
	})

	t.Run("invalid chunk size", func(t *testing.T) {
		chunkedData := "XYZ\r\nHello\r\n"

		reader := NewChunkedReader(strings.NewReader(chunkedData))
		result := make([]byte, 10)
		_, err := reader.Read(result)

		if err == nil {
			t.Error("Expected error for invalid chunk size")
		}
	})
}

func TestContentLengthReader(t *testing.T) {
	t.Run("read with content length", func(t *testing.T) {
		data := "Hello, World!"
		reader := NewContentLengthReader(strings.NewReader(data), 5)

		result := make([]byte, 10)
		n, err := reader.Read(result)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "Hello"
		if string(result[:n]) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result[:n]))
		}

		if reader.Remaining() != 0 {
			t.Errorf("Expected 0 remaining, got %d", reader.Remaining())
		}
	})

	t.Run("read beyond content length", func(t *testing.T) {
		data := "Hello"
		reader := NewContentLengthReader(strings.NewReader(data), 3)

		result := make([]byte, 10)
		n, err := reader.Read(result)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "Hel"
		if string(result[:n]) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result[:n]))
		}
	})
}

func TestHTTPMessageBuilder(t *testing.T) {
	builder := NewHTTPMessageBuilder()

	t.Run("build request", func(t *testing.T) {
		req := pkghttp.NewRequest(pkghttp.MethodGet, "/test", pkghttp.Version11)
		req.SetHeader("Host", "example.com")

		data, err := builder.BuildRequest(req)
		if err != nil {
			t.Fatalf("BuildRequest failed: %v", err)
		}

		expected := "GET /test HTTP/1.1\r\nHost: example.com\r\n\r\n"
		if string(data) != expected {
			t.Errorf("Expected:\n%s\nGot:\n%s", expected, string(data))
		}
	})

	t.Run("build response", func(t *testing.T) {
		resp := pkghttp.NewTextResponse(pkghttp.StatusOK, pkghttp.Version11, "Hello")

		data, err := builder.BuildResponse(resp)
		if err != nil {
			t.Fatalf("BuildResponse failed: %v", err)
		}

		if !bytes.Contains(data, []byte("HTTP/1.1 200 OK")) {
			t.Error("Response should contain status line")
		}

		if !bytes.Contains(data, []byte("Hello")) {
			t.Error("Response should contain body")
		}
	})
}

func TestParseChunkSize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "simple hex",
			input:    "a",
			expected: 10,
			wantErr:  false,
		},
		{
			name:     "uppercase hex",
			input:    "A",
			expected: 10,
			wantErr:  false,
		},
		{
			name:     "decimal",
			input:    "10",
			expected: 16,
			wantErr:  false,
		},
		{
			name:     "with extension",
			input:    "a;name=value",
			expected: 10,
			wantErr:  false,
		},
		{
			name:     "invalid character",
			input:    "xyz",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseChunkSize(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestValidationFunctions(t *testing.T) {
	t.Run("isValidMethod", func(t *testing.T) {
		tests := []struct {
			method pkghttp.Method
			valid  bool
		}{
			{pkghttp.MethodGet, true},
			{pkghttp.MethodPost, true},
			{pkghttp.MethodPut, true},
			{pkghttp.MethodDelete, true},
			{pkghttp.MethodHead, true},
			{pkghttp.MethodOptions, true},
			{pkghttp.MethodPatch, true},
			{"INVALID", false},
		}

		for _, tt := range tests {
			if isValidMethod(tt.method) != tt.valid {
				t.Errorf("isValidMethod(%s) = %t, want %t", tt.method, !tt.valid, tt.valid)
			}
		}
	})

	t.Run("isValidPath", func(t *testing.T) {
		tests := []struct {
			path  string
			valid bool
		}{
			{"/", true},
			{"/hello", true},
			{"/api/v1/users", true},
			{"/path?query=value", true},
			{"", false},
			{"hello", false},
			{"/path\x00", false},
		}

		for _, tt := range tests {
			if isValidPath(tt.path) != tt.valid {
				t.Errorf("isValidPath(%s) = %t, want %t", tt.path, !tt.valid, tt.valid)
			}
		}
	})

	t.Run("isValidVersion", func(t *testing.T) {
		tests := []struct {
			version pkghttp.Version
			valid   bool
		}{
			{pkghttp.Version10, true},
			{pkghttp.Version11, true},
			{"HTTP/2.0", false},
			{"", false},
		}

		for _, tt := range tests {
			if isValidVersion(tt.version) != tt.valid {
				t.Errorf("isValidVersion(%s) = %t, want %t", tt.version, !tt.valid, tt.valid)
			}
		}
	})

	t.Run("isValidHeaderName", func(t *testing.T) {
		tests := []struct {
			name  string
			valid bool
		}{
			{"Host", true},
			{"Content-Type", true},
			{"X-Custom-Header", true},
			{"", false},
			{"Header with spaces", false},
			{"Header:with:colons", false},
		}

		for _, tt := range tests {
			if isValidHeaderName(tt.name) != tt.valid {
				t.Errorf("isValidHeaderName(%s) = %t, want %t", tt.name, !tt.valid, tt.valid)
			}
		}
	})
}