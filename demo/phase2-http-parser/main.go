package main

import (
	"fmt"
	"strings"

	"github.com/ganyariya/tinyserver/internal/http"
	pkghttp "github.com/ganyariya/tinyserver/pkg/http"
)

func main() {
	fmt.Println("HTTP Request Parser & Analyzer Demo")
	fmt.Println("====================================")
	fmt.Println()

	// Create HTTP parser
	parser := http.NewParser()

	// Sample HTTP requests for demonstration
	samples := []string{
		"GET /api/users?id=123 HTTP/1.1\r\nHost: example.com\r\nUser-Agent: TinyClient/1.0\r\nAccept: application/json\r\n\r\n",
		"POST /api/login HTTP/1.1\r\nHost: api.example.com\r\nContent-Type: application/json\r\nContent-Length: 40\r\nAuthorization: Bearer token123\r\n\r\n{\"username\":\"user\",\"password\":\"pass123\"}",
		"PUT /api/users/456 HTTP/1.1\r\nHost: api.example.com\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nUpdated data!",
		"DELETE /api/users/789 HTTP/1.1\r\nHost: api.example.com\r\nAuthorization: Bearer token456\r\n\r\n",
		"GET /invalid request line", // Invalid request for error demonstration
	}

	// Parse and analyze each sample
	for i, rawRequest := range samples {
		fmt.Printf("=== Sample Request %d ===\n", i+1)
		fmt.Println("Raw HTTP Request:")
		fmt.Println(strings.ReplaceAll(rawRequest, "\r\n", "\\r\\n\n"))
		fmt.Println()

		// Parse the request
		req, err := parser.ParseBytes([]byte(rawRequest))
		if err != nil {
			fmt.Printf("❌ Parse Error: %v\n", err)
			fmt.Println()
			continue
		}

		// Validate the request
		if validationErr := parser.Validate(req); validationErr != nil {
			fmt.Printf("❌ Validation Error: %v\n", validationErr)
			fmt.Println()
			continue
		}

		// Display parsed components
		fmt.Println("✅ Parse Result:")
		fmt.Printf("  Method: %s\n", req.Method())
		fmt.Printf("  Path: %s\n", req.Path())
		fmt.Printf("  Version: %s\n", req.Version())

		// Display query parameters
		queryParams := req.QueryParams()
		if len(queryParams) > 0 {
			fmt.Println("  Query Parameters:")
			for key, value := range queryParams {
				fmt.Printf("    %s: %s\n", key, value)
			}
		}

		// Display headers
		headers := req.Headers()
		if len(headers) > 0 {
			fmt.Println("  Headers:")
			for name, values := range headers {
				for _, value := range values {
					fmt.Printf("    %s: %s\n", name, value)
				}
			}
		}

		// Display content length
		contentLength := req.ContentLength()
		fmt.Printf("  Content Length: %d\n", contentLength)

		// Display body if present
		if contentLength > 0 {
			body := make([]byte, contentLength)
			n, readErr := req.Body().Read(body)
			if readErr == nil && n > 0 {
				fmt.Printf("  Body: %s\n", string(body[:n]))
			}
		}

		// Generate response for demonstration
		fmt.Println()
		fmt.Println("Generated Response:")
		response := generateResponse(req)
		fmt.Println(response)
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println()
	}

	// Demonstrate response parsing
	fmt.Println("=== HTTP Response Parser Demo ===")

	responseParser := http.NewResponseParser()
	sampleResponse := "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 50\r\nServer: TinyServer/1.0\r\n\r\n{\"status\":\"success\",\"message\":\"Request processed\"}"

	fmt.Println("Raw HTTP Response:")
	fmt.Println(strings.ReplaceAll(sampleResponse, "\r\n", "\\r\\n\n"))
	fmt.Println()

	resp, err := responseParser.ParseResponseBytes([]byte(sampleResponse))
	if err != nil {
		fmt.Printf("❌ Response Parse Error: %v\n", err)
	} else {
		fmt.Println("✅ Response Parse Result:")
		fmt.Printf("  Status Code: %d\n", resp.StatusCode())
		fmt.Printf("  Version: %s\n", resp.Version())
		fmt.Printf("  Content Length: %d\n", resp.ContentLength())

		respHeaders := resp.Headers()
		if len(respHeaders) > 0 {
			fmt.Println("  Headers:")
			for name, values := range respHeaders {
				for _, value := range values {
					fmt.Printf("    %s: %s\n", name, value)
				}
			}
		}

		if resp.ContentLength() > 0 {
			body := make([]byte, resp.ContentLength())
			n, readErr := resp.Body().Read(body)
			if readErr == nil && n > 0 {
				fmt.Printf("  Body: %s\n", string(body[:n]))
			}
		}
	}

	fmt.Println()
	fmt.Println("🎉 HTTP Parser Demo Complete!")
	fmt.Println("This demo showcased:")
	fmt.Println("  ✓ HTTP request parsing and validation")
	fmt.Println("  ✓ Query parameter extraction")
	fmt.Println("  ✓ Header parsing and display")
	fmt.Println("  ✓ Request body handling")
	fmt.Println("  ✓ HTTP response parsing")
	fmt.Println("  ✓ Error handling for malformed requests")
}

// generateResponse creates a sample HTTP response for demonstration
func generateResponse(req pkghttp.Request) string {
	method := string(req.Method())
	path := req.Path()

	switch method {
	case "GET":
		if strings.Contains(path, "/api/users") {
			body := `{"id":123,"name":"John Doe","email":"john@example.com"}`
			return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
				pkghttp.StatusOK, pkghttp.StatusText(pkghttp.StatusOK),
				pkghttp.HeaderContentType, pkghttp.MimeTypeJSON,
				pkghttp.HeaderContentLength, len(body),
				body)
		}
		body := "GET request processed"
		return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
			pkghttp.StatusOK, pkghttp.StatusText(pkghttp.StatusOK),
			pkghttp.HeaderContentType, pkghttp.MimeTypeTextPlain,
			pkghttp.HeaderContentLength, len(body),
			body)

	case "POST":
		if strings.Contains(path, "/login") {
			body := `{"status":"success","token":"abc123xyz"}`
			return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
				pkghttp.StatusOK, pkghttp.StatusText(pkghttp.StatusOK),
				pkghttp.HeaderContentType, pkghttp.MimeTypeJSON,
				pkghttp.HeaderContentLength, len(body),
				body)
		}
		body := "POST request processed"
		return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
			pkghttp.StatusCreated, pkghttp.StatusText(pkghttp.StatusCreated),
			pkghttp.HeaderContentType, pkghttp.MimeTypeTextPlain,
			pkghttp.HeaderContentLength, len(body),
			body)

	case "PUT":
		body := "PUT request processed"
		return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
			pkghttp.StatusOK, pkghttp.StatusText(pkghttp.StatusOK),
			pkghttp.HeaderContentType, pkghttp.MimeTypeTextPlain,
			pkghttp.HeaderContentLength, len(body),
			body)

	case "DELETE":
		return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: 0\r\n\r\n",
			pkghttp.StatusNoContent, pkghttp.StatusText(pkghttp.StatusNoContent),
			pkghttp.HeaderContentLength)

	default:
		body := "Method not allowed"
		return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s: %s\r\n%s: %d\r\n\r\n%s",
			pkghttp.StatusMethodNotAllowed, pkghttp.StatusText(pkghttp.StatusMethodNotAllowed),
			pkghttp.HeaderContentType, pkghttp.MimeTypeTextPlain,
			pkghttp.HeaderContentLength, len(body),
			body)
	}
}
