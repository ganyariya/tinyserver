package http

import (
	"strings"
	"testing"

	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

func TestWriteRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  pkghttp.Request
		expected string
	}{
		{
			name: "simple GET request",
			request: func() pkghttp.Request {
				req := pkghttp.NewRequest(pkghttp.MethodGet, "/hello", pkghttp.Version11)
				req.SetHeader("Host", "example.com")
				req.SetHeader("User-Agent", "TinyClient/1.0")
				return req
			}(),
			expected: "GET /hello HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"User-Agent: TinyClient/1.0\r\n" +
				"\r\n",
		},
		{
			name: "POST request with body",
			request: func() pkghttp.Request {
				req := pkghttp.NewRequest(pkghttp.MethodPost, "/api/data", pkghttp.Version11)
				req.SetHeader("Host", "example.com")
				req.SetHeader("Content-Type", "application/json")
				req.SetHeader("Content-Length", "13")
				req.SetBody(strings.NewReader("{\"test\": true}"))
				return req
			}(),
			expected: "POST /api/data HTTP/1.1\r\n" +
				"Host: example.com\r\n" +
				"Content-Type: application/json\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"{\"test\": true}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := WriteRequest(&buf, tt.request)

			if err != nil {
				t.Fatalf("WriteRequest failed: %v", err)
			}

			result := buf.String()

			// Check request line
			lines := strings.Split(result, "\r\n")
			if len(lines) < 2 {
				t.Fatalf("Request has too few lines: %d", len(lines))
			}

			expectedLines := strings.Split(tt.expected, "\r\n")
			requestLine := lines[0]
			expectedRequestLine := expectedLines[0]
			if requestLine != expectedRequestLine {
				t.Errorf("Request line mismatch:\nExpected: %s\nGot: %s", expectedRequestLine, requestLine)
			}

			// Check that all expected headers are present
			headerLines := make(map[string]string)
			for i := 1; i < len(lines) && lines[i] != ""; i++ {
				parts := strings.SplitN(lines[i], ": ", 2)
				if len(parts) == 2 {
					headerLines[parts[0]] = parts[1]
				}
			}

			expectedHeaders := make(map[string]string)
			for i := 1; i < len(expectedLines) && expectedLines[i] != ""; i++ {
				parts := strings.SplitN(expectedLines[i], ": ", 2)
				if len(parts) == 2 {
					expectedHeaders[parts[0]] = parts[1]
				}
			}

			for expectedHeader, expectedValue := range expectedHeaders {
				if actualValue, exists := headerLines[expectedHeader]; !exists {
					t.Errorf("Missing header: %s", expectedHeader)
				} else if actualValue != expectedValue {
					t.Errorf("Header value mismatch for %s:\nExpected: %s\nGot: %s", expectedHeader, expectedValue, actualValue)
				}
			}

			// Check body if present
			bodyStart := strings.Index(result, "\r\n\r\n")
			expectedBodyStart := strings.Index(tt.expected, "\r\n\r\n")
			if bodyStart != -1 && expectedBodyStart != -1 {
				actualBody := result[bodyStart+4:]
				expectedBody := tt.expected[expectedBodyStart+4:]
				if actualBody != expectedBody {
					t.Errorf("Body mismatch:\nExpected: %s\nGot: %s", expectedBody, actualBody)
				}
			}
		})
	}
}

func TestFormatRequest(t *testing.T) {
	req := pkghttp.NewRequest(pkghttp.MethodGet, "/hello?name=world", pkghttp.Version11)
	req.SetHeader("Host", "example.com")
	req.SetHeader("User-Agent", "TinyClient/1.0")

	formatted := FormatRequest(req)

	expectedContains := []string{
		"GET /hello?name=world HTTP/1.1",
		"Host: example.com",
		"User-Agent: TinyClient/1.0",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(formatted, expected) {
			t.Errorf("Formatted request should contain: %s\nGot:\n%s", expected, formatted)
		}
	}
}

func TestParseRequestLine(t *testing.T) {
	tests := []struct {
		name        string
		requestLine string
		wantMethod  pkghttp.Method
		wantPath    string
		wantVersion pkghttp.Version
		wantErr     bool
	}{
		{
			name:        "valid GET request",
			requestLine: "GET /hello HTTP/1.1",
			wantMethod:  pkghttp.MethodGet,
			wantPath:    "/hello",
			wantVersion: pkghttp.Version11,
			wantErr:     false,
		},
		{
			name:        "valid POST request",
			requestLine: "POST /api/data HTTP/1.0",
			wantMethod:  pkghttp.MethodPost,
			wantPath:    "/api/data",
			wantVersion: pkghttp.Version10,
			wantErr:     false,
		},
		{
			name:        "empty line",
			requestLine: "",
			wantErr:     true,
		},
		{
			name:        "too few parts",
			requestLine: "GET /hello",
			wantErr:     true,
		},
		{
			name:        "invalid method",
			requestLine: "INVALID /hello HTTP/1.1",
			wantErr:     true,
		},
		{
			name:        "invalid path",
			requestLine: "GET hello HTTP/1.1",
			wantErr:     true,
		},
		{
			name:        "invalid version",
			requestLine: "GET /hello HTTP/2.0",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, path, version, err := parseRequestLine(tt.requestLine)

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

			if method != tt.wantMethod {
				t.Errorf("Expected method %s, got %s", tt.wantMethod, method)
			}

			if path != tt.wantPath {
				t.Errorf("Expected path %s, got %s", tt.wantPath, path)
			}

			if version != tt.wantVersion {
				t.Errorf("Expected version %s, got %s", tt.wantVersion, version)
			}
		})
	}
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name       string
		headerLine string
		wantName   string
		wantValue  string
		wantErr    bool
	}{
		{
			name:       "simple header",
			headerLine: "Host: example.com",
			wantName:   "Host",
			wantValue:  "example.com",
			wantErr:    false,
		},
		{
			name:       "header with spaces",
			headerLine: "Content-Type:  application/json  ",
			wantName:   "Content-Type",
			wantValue:  "application/json",
			wantErr:    false,
		},
		{
			name:       "header with colon in value",
			headerLine: "Authorization: Bearer token:with:colons",
			wantName:   "Authorization",
			wantValue:  "Bearer token:with:colons",
			wantErr:    false,
		},
		{
			name:       "no colon",
			headerLine: "Invalid header line",
			wantErr:    true,
		},
		{
			name:       "empty name",
			headerLine: ": value",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, value, err := parseHeader(tt.headerLine)

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

			if name != tt.wantName {
				t.Errorf("Expected name %s, got %s", tt.wantName, name)
			}

			if value != tt.wantValue {
				t.Errorf("Expected value %s, got %s", tt.wantValue, value)
			}
		})
	}
}

func TestNewRequestFromRaw(t *testing.T) {
	rawData := []byte("GET /hello HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"User-Agent: TinyClient/1.0\r\n" +
		"\r\n")

	req, err := NewRequestFromRaw(rawData, nil)
	if err != nil {
		t.Fatalf("NewRequestFromRaw failed: %v", err)
	}

	if req.Method() != pkghttp.MethodGet {
		t.Errorf("Expected GET, got %s", req.Method())
	}

	if req.Path() != "/hello" {
		t.Errorf("Expected /hello, got %s", req.Path())
	}

	if req.Version() != pkghttp.Version11 {
		t.Errorf("Expected HTTP/1.1, got %s", req.Version())
	}

	if req.GetHeader("Host") != "example.com" {
		t.Errorf("Expected Host header value example.com, got %s", req.GetHeader("Host"))
	}
}

func TestParseRequestWithBody(t *testing.T) {
	rawData := "POST /api/data HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Type: application/json\r\n" +
		"Content-Length: 14\r\n" +
		"\r\n" +
		"{\"test\": true}"

	reader := strings.NewReader(rawData)
	req, err := ParseRequest(reader, nil)

	if err != nil {
		t.Fatalf("ParseRequest failed: %v", err)
	}

	if req.Method() != pkghttp.MethodPost {
		t.Errorf("Expected POST, got %s", req.Method())
	}

	if req.ContentLength() != 14 {
		t.Errorf("Expected content length 14, got %d", req.ContentLength())
	}

	if req.GetHeader("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", req.GetHeader("Content-Type"))
	}
}

func TestParseRequestErrors(t *testing.T) {
	tests := []struct {
		name    string
		rawData string
	}{
		{
			name:    "empty request",
			rawData: "",
		},
		{
			name:    "invalid request line",
			rawData: "INVALID\r\n\r\n",
		},
		{
			name: "invalid header",
			rawData: "GET /hello HTTP/1.1\r\n" +
				"Invalid header line\r\n" +
				"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.rawData)
			_, err := ParseRequest(reader, nil)

			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}
