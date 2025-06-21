# Phase 2: HTTP Parser & Analyzer Demo

This demo showcases the HTTP request/response parsing capabilities implemented in Phase 2 of the TinyServer project.

## 🎯 Demo Overview

The HTTP Parser Demo demonstrates:
- **HTTP Request Parsing**: Parse various HTTP methods (GET, POST, PUT, DELETE)
- **Header Analysis**: Extract and display HTTP headers
- **Query Parameter Extraction**: Parse URL query parameters
- **Request Body Handling**: Process request bodies with content-length
- **Response Generation**: Create appropriate HTTP responses
- **Error Handling**: Handle malformed requests gracefully
- **Validation**: Validate parsed HTTP messages

## 🚀 Running the Demo

### Prerequisites
- Go 1.19 or later
- TinyServer project built (`go build ./...`)

### Quick Start
```bash
# From the project root directory
go run demo/phase2-http-parser/main.go
```

### Alternative: Using Make (if available)
```bash
make demo-phase2
```

### Alternative: Using Demo Script
```bash
./scripts/demo/run-phase2.sh
```

## 🔍 What You'll See

The demo processes several sample HTTP requests and shows:

### Sample Request Analysis
```
=== Sample Request 1 ===
Raw HTTP Request:
GET /api/users?id=123 HTTP/1.1\r\n
Host: example.com\r\n
User-Agent: TinyClient/1.0\r\n
Accept: application/json\r\n\r\n

✅ Parse Result:
  Method: GET
  Path: /api/users
  Version: HTTP/1.1
  Query Parameters:
    id: 123
  Headers:
    Host: example.com
    User-Agent: TinyClient/1.0
    Accept: application/json
  Content Length: 0

Generated Response:
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 52

{"id":123,"name":"John Doe","email":"john@example.com"}
```

### Features Demonstrated

1. **HTTP Methods**: GET, POST, PUT, DELETE requests
2. **Query Parameters**: URL parameter parsing and extraction
3. **Headers**: Complete header parsing and display
4. **Request Bodies**: JSON and text body handling
5. **Response Generation**: Appropriate response creation
6. **Error Handling**: Malformed request detection
7. **Validation**: HTTP message validation

## 📋 Technical Details

### HTTP Parser Capabilities
- **Request Line Parsing**: Method, path, and version extraction
- **Header Processing**: Multi-value header support
- **Body Reading**: Content-length based body reading
- **Query String Parsing**: URL parameter extraction
- **Protocol Validation**: HTTP/1.1 compliance checking

### Supported HTTP Methods
- `GET` - Retrieve data
- `POST` - Create new resources
- `PUT` - Update existing resources
- `DELETE` - Remove resources
- `HEAD` - Retrieve headers only
- `OPTIONS` - Query allowed methods
- `PATCH` - Partial updates

### Header Handling
- Standard HTTP headers (Host, User-Agent, Content-Type, etc.)
- Authorization headers
- Custom headers
- Multi-value headers
- Case-insensitive header names

## 🎓 Learning Objectives

This demo helps you understand:

1. **HTTP Protocol Structure**: How HTTP messages are formatted
2. **Parsing Techniques**: String processing and protocol parsing
3. **Error Handling**: Graceful handling of malformed requests
4. **Data Validation**: Ensuring protocol compliance
5. **Response Generation**: Creating proper HTTP responses

## 🔧 Customization

You can modify the demo to:
- Add more sample requests
- Test different HTTP methods
- Experiment with various header combinations
- Add custom validation rules
- Extend response generation logic

## 🧪 Testing Your Understanding

Try these exercises:
1. Add a new HTTP method (e.g., PATCH)
2. Create a request with invalid headers
3. Test with different Content-Type values
4. Add custom header validation
5. Implement chunked transfer encoding parsing

## 📊 Performance Notes

The parser demonstrates:
- Efficient string processing
- Memory-conscious header parsing
- Streaming body reading
- Timeout handling capabilities

## 🎉 Next Steps

After running this demo, you'll be ready for:
- **Phase 3**: Building a simple HTTP server
- **Phase 4**: Implementing an HTTP client
- **Phase 5**: Creating a full-stack web application

## 🐛 Troubleshooting

If you encounter issues:

1. **Build Errors**: Ensure all dependencies are installed with `go mod tidy`
2. **Import Errors**: Verify you're running from the project root directory
3. **Parse Errors**: Check that sample requests are properly formatted
4. **Missing Functions**: Ensure Phase 2 implementation is complete

## 📚 Related Files

- `internal/http/parser.go` - Main HTTP parser implementation
- `internal/http/request.go` - HTTP request parsing logic
- `internal/http/response.go` - HTTP response generation
- `pkg/http/interfaces.go` - HTTP interface definitions
- `pkg/http/constants.go` - HTTP constants and status codes

---

This demo represents a crucial step in understanding HTTP protocol implementation. The parsing capabilities demonstrated here form the foundation for building web servers and clients in the subsequent phases.