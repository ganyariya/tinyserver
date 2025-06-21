package http

import (
	"strings"
	"testing"

	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		name     string
		response pkghttp.Response
		expected string
	}{
		{
			name: "simple OK response",
			response: func() pkghttp.Response {
				resp := pkghttp.NewTextResponse(pkghttp.StatusOK, pkghttp.Version11, "Hello, World!")
				resp.SetHeader("Server", "TinyServer/1.0")
				return resp
			}(),
			expected: "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: 13\r\n" +
				"Server: TinyServer/1.0\r\n" +
				"\r\n" +
				"Hello, World!",
		},
		{
			name: "not found response",
			response: func() pkghttp.Response {
				resp := pkghttp.NewHTMLResponse(pkghttp.StatusNotFound, pkghttp.Version11, "<h1>Not Found</h1>")
				return resp
			}(),
			expected: "HTTP/1.1 404 Not Found\r\n" +
				"Content-Type: text/html\r\n" +
				"Content-Length: 18\r\n" +
				"\r\n" +
				"<h1>Not Found</h1>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := WriteResponse(&buf, tt.response)

			if err != nil {
				t.Fatalf("WriteResponse failed: %v", err)
			}

			result := buf.String()
			
			// Check status line
			lines := strings.Split(result, "\r\n")
			if len(lines) < 3 {
				t.Fatalf("Response has too few lines: %d", len(lines))
			}
			
			expectedLines := strings.Split(tt.expected, "\r\n")
			statusLine := lines[0]
			expectedStatusLine := expectedLines[0]
			if statusLine != expectedStatusLine {
				t.Errorf("Status line mismatch:\nExpected: %s\nGot: %s", expectedStatusLine, statusLine)
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
			
			// Check body
			bodyStart := strings.Index(result, "\r\n\r\n")
			expectedBodyStart := strings.Index(tt.expected, "\r\n\r\n")
			if bodyStart == -1 || expectedBodyStart == -1 {
				t.Fatal("Could not find body separator")
			}
			
			actualBody := result[bodyStart+4:]
			expectedBody := tt.expected[expectedBodyStart+4:]
			if actualBody != expectedBody {
				t.Errorf("Body mismatch:\nExpected: %s\nGot: %s", expectedBody, actualBody)
			}
		})
	}
}

func TestFormatResponse(t *testing.T) {
	resp := pkghttp.NewTextResponse(pkghttp.StatusOK, pkghttp.Version11, "Hello")
	resp.SetHeader("Server", "TinyServer/1.0")

	formatted := FormatResponse(resp)

	expectedContains := []string{
		"HTTP/1.1 200 OK",
		"Server: TinyServer/1.0",
		"Success: true",
		"Error: false",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(formatted, expected) {
			t.Errorf("Formatted response should contain: %s\nGot:\n%s", expected, formatted)
		}
	}
}

func TestParseStatusLine(t *testing.T) {
	tests := []struct {
		name           string
		statusLine     string
		wantVersion    pkghttp.Version
		wantStatusCode pkghttp.StatusCode
		wantErr        bool
	}{
		{
			name:           "valid OK response",
			statusLine:     "HTTP/1.1 200 OK",
			wantVersion:    pkghttp.Version11,
			wantStatusCode: pkghttp.StatusOK,
			wantErr:        false,
		},
		{
			name:           "valid not found response",
			statusLine:     "HTTP/1.0 404 Not Found",
			wantVersion:    pkghttp.Version10,
			wantStatusCode: pkghttp.StatusNotFound,
			wantErr:        false,
		},
		{
			name:        "empty line",
			statusLine:  "",
			wantErr:     true,
		},
		{
			name:        "too few parts",
			statusLine:  "HTTP/1.1",
			wantErr:     true,
		},
		{
			name:        "invalid version",
			statusLine:  "HTTP/2.0 200 OK",
			wantErr:     true,
		},
		{
			name:        "invalid status code",
			statusLine:  "HTTP/1.1 999 Unknown",
			wantErr:     true,
		},
		{
			name:        "non-numeric status code",
			statusLine:  "HTTP/1.1 ABC OK",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, statusCode, err := parseStatusLine(tt.statusLine)

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

			if version != tt.wantVersion {
				t.Errorf("Expected version %s, got %s", tt.wantVersion, version)
			}

			if statusCode != tt.wantStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.wantStatusCode, statusCode)
			}
		})
	}
}

func TestBuildErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode pkghttp.StatusCode
		message    string
	}{
		{
			name:       "not found error",
			statusCode: pkghttp.StatusNotFound,
			message:    "The requested resource was not found",
		},
		{
			name:       "internal server error",
			statusCode: pkghttp.StatusInternalServerError,
			message:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := BuildErrorResponse(tt.statusCode, tt.message)

			if resp.StatusCode() != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, resp.StatusCode())
			}

			if resp.GetHeader("Content-Type") != pkghttp.MimeTypeTextHTML {
				t.Errorf("Expected Content-Type %s, got %s", pkghttp.MimeTypeTextHTML, resp.GetHeader("Content-Type"))
			}

			// Check that response contains HTML structure
			var buf strings.Builder
			resp.WriteTo(&buf)
			responseText := buf.String()

			if !strings.Contains(responseText, "<html>") {
				t.Error("Error response should contain HTML")
			}

			if !strings.Contains(responseText, "TinyServer") {
				t.Error("Error response should contain server name")
			}
		})
	}
}

func TestBuildJSONErrorResponse(t *testing.T) {
	resp := BuildJSONErrorResponse(pkghttp.StatusBadRequest, "Invalid input")

	if resp.StatusCode() != pkghttp.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", pkghttp.StatusBadRequest, resp.StatusCode())
	}

	if resp.GetHeader("Content-Type") != pkghttp.MimeTypeJSON {
		t.Errorf("Expected Content-Type %s, got %s", pkghttp.MimeTypeJSON, resp.GetHeader("Content-Type"))
	}

	var buf strings.Builder
	resp.WriteTo(&buf)
	responseText := buf.String()

	if !strings.Contains(responseText, `"error"`) {
		t.Error("JSON error response should contain error field")
	}

	if !strings.Contains(responseText, "Invalid input") {
		t.Error("JSON error response should contain error message")
	}
}

func TestBuildRedirectResponse(t *testing.T) {
	location := "https://example.com/new-location"
	resp := BuildRedirectResponse(pkghttp.StatusMovedPermanently, location)

	if resp.StatusCode() != pkghttp.StatusMovedPermanently {
		t.Errorf("Expected status code %d, got %d", pkghttp.StatusMovedPermanently, resp.StatusCode())
	}

	if resp.GetHeader("Location") != location {
		t.Errorf("Expected Location header %s, got %s", location, resp.GetHeader("Location"))
	}

	if resp.GetHeader("Content-Type") != pkghttp.MimeTypeTextHTML {
		t.Errorf("Expected Content-Type %s, got %s", pkghttp.MimeTypeTextHTML, resp.GetHeader("Content-Type"))
	}

	var buf strings.Builder
	resp.WriteTo(&buf)
	responseText := buf.String()

	if !strings.Contains(responseText, location) {
		t.Error("Redirect response should contain redirect URL in body")
	}
}

func TestSetCommonHeaders(t *testing.T) {
	resp := pkghttp.NewResponse(pkghttp.StatusOK, pkghttp.Version11)
	SetCommonHeaders(resp)

	if resp.GetHeader("Server") != "TinyServer/1.0" {
		t.Errorf("Expected Server header TinyServer/1.0, got %s", resp.GetHeader("Server"))
	}

	if resp.GetHeader("Connection") != "close" {
		t.Errorf("Expected Connection header close, got %s", resp.GetHeader("Connection"))
	}

	if resp.GetHeader("Date") == "" {
		t.Error("Date header should be set")
	}
}

func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name     string
		response pkghttp.Response
		wantErr  bool
	}{
		{
			name:     "valid response",
			response: pkghttp.NewTextResponse(pkghttp.StatusOK, pkghttp.Version11, "Hello"),
			wantErr:  false,
		},
		{
			name: "invalid status code",
			response: func() pkghttp.Response {
				resp := pkghttp.NewResponse(pkghttp.StatusCode(999), pkghttp.Version11)
				return resp
			}(),
			wantErr: true,
		},
		{
			name: "invalid version",
			response: func() pkghttp.Response {
				resp := pkghttp.NewResponse(pkghttp.StatusOK, "HTTP/2.0")
				return resp
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateResponse(tt.response)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestNewResponseFromRaw(t *testing.T) {
	rawData := []byte("HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 5\r\n" +
		"\r\n" +
		"Hello")

	resp, err := NewResponseFromRaw(rawData)
	if err != nil {
		t.Fatalf("NewResponseFromRaw failed: %v", err)
	}

	if resp.StatusCode() != pkghttp.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode())
	}

	if resp.Version() != pkghttp.Version11 {
		t.Errorf("Expected HTTP/1.1, got %s", resp.Version())
	}

	if resp.GetHeader("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type text/plain, got %s", resp.GetHeader("Content-Type"))
	}

	if resp.ContentLength() != 5 {
		t.Errorf("Expected content length 5, got %d", resp.ContentLength())
	}
}

func TestParseResponseErrors(t *testing.T) {
	tests := []struct {
		name    string
		rawData string
	}{
		{
			name:    "empty response",
			rawData: "",
		},
		{
			name:    "invalid status line",
			rawData: "INVALID\r\n\r\n",
		},
		{
			name: "invalid header",
			rawData: "HTTP/1.1 200 OK\r\n" +
				"Invalid header line\r\n" +
				"\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.rawData)
			_, err := ParseResponse(reader)

			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}