package http

import "time"

// Common HTTP status codes
const (
	// 1xx Informational
	StatusContinue           StatusCode = 100
	StatusSwitchingProtocols StatusCode = 101

	// 2xx Success
	StatusOK                   StatusCode = 200
	StatusCreated              StatusCode = 201
	StatusAccepted             StatusCode = 202
	StatusNonAuthoritativeInfo StatusCode = 203
	StatusNoContent            StatusCode = 204
	StatusResetContent         StatusCode = 205
	StatusPartialContent       StatusCode = 206

	// 3xx Redirection
	StatusMultipleChoices   StatusCode = 300
	StatusMovedPermanently  StatusCode = 301
	StatusFound             StatusCode = 302
	StatusSeeOther          StatusCode = 303
	StatusNotModified       StatusCode = 304
	StatusUseProxy          StatusCode = 305
	StatusTemporaryRedirect StatusCode = 307
	StatusPermanentRedirect StatusCode = 308

	// 4xx Client Error
	StatusBadRequest                   StatusCode = 400
	StatusUnauthorized                 StatusCode = 401
	StatusPaymentRequired              StatusCode = 402
	StatusForbidden                    StatusCode = 403
	StatusNotFound                     StatusCode = 404
	StatusMethodNotAllowed             StatusCode = 405
	StatusNotAcceptable                StatusCode = 406
	StatusProxyAuthRequired            StatusCode = 407
	StatusRequestTimeout               StatusCode = 408
	StatusConflict                     StatusCode = 409
	StatusGone                         StatusCode = 410
	StatusLengthRequired               StatusCode = 411
	StatusPreconditionFailed           StatusCode = 412
	StatusRequestEntityTooLarge        StatusCode = 413
	StatusRequestURITooLong            StatusCode = 414
	StatusUnsupportedMediaType         StatusCode = 415
	StatusRequestedRangeNotSatisfiable StatusCode = 416
	StatusExpectationFailed            StatusCode = 417
	StatusTeapot                       StatusCode = 418
	StatusMisdirectedRequest           StatusCode = 421
	StatusUnprocessableEntity          StatusCode = 422
	StatusLocked                       StatusCode = 423
	StatusFailedDependency             StatusCode = 424
	StatusTooEarly                     StatusCode = 425
	StatusUpgradeRequired              StatusCode = 426
	StatusPreconditionRequired         StatusCode = 428
	StatusTooManyRequests              StatusCode = 429
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431
	StatusUnavailableForLegalReasons   StatusCode = 451

	// 5xx Server Error
	StatusInternalServerError           StatusCode = 500
	StatusNotImplemented                StatusCode = 501
	StatusBadGateway                    StatusCode = 502
	StatusServiceUnavailable            StatusCode = 503
	StatusGatewayTimeout                StatusCode = 504
	StatusHTTPVersionNotSupported       StatusCode = 505
	StatusVariantAlsoNegotiates         StatusCode = 506
	StatusInsufficientStorage           StatusCode = 507
	StatusLoopDetected                  StatusCode = 508
	StatusNotExtended                   StatusCode = 510
	StatusNetworkAuthenticationRequired StatusCode = 511
)

// Common HTTP headers
const (
	HeaderAccept                          = "Accept"
	HeaderAcceptCharset                   = "Accept-Charset"
	HeaderAcceptEncoding                  = "Accept-Encoding"
	HeaderAcceptLanguage                  = "Accept-Language"
	HeaderAcceptRanges                    = "Accept-Ranges"
	HeaderAge                             = "Age"
	HeaderAllow                           = "Allow"
	HeaderAuthorization                   = "Authorization"
	HeaderCacheControl                    = "Cache-Control"
	HeaderConnection                      = "Connection"
	HeaderContentDisposition              = "Content-Disposition"
	HeaderContentEncoding                 = "Content-Encoding"
	HeaderContentLanguage                 = "Content-Language"
	HeaderContentLength                   = "Content-Length"
	HeaderContentLocation                 = "Content-Location"
	HeaderContentRange                    = "Content-Range"
	HeaderContentType                     = "Content-Type"
	HeaderDate                            = "Date"
	HeaderETag                            = "ETag"
	HeaderExpect                          = "Expect"
	HeaderExpires                         = "Expires"
	HeaderFrom                            = "From"
	HeaderHost                            = "Host"
	HeaderIfMatch                         = "If-Match"
	HeaderIfModifiedSince                 = "If-Modified-Since"
	HeaderIfNoneMatch                     = "If-None-Match"
	HeaderIfRange                         = "If-Range"
	HeaderIfUnmodifiedSince               = "If-Unmodified-Since"
	HeaderLastModified                    = "Last-Modified"
	HeaderLocation                        = "Location"
	HeaderMaxForwards                     = "Max-Forwards"
	HeaderPragma                          = "Pragma"
	HeaderProxyAuthenticate               = "Proxy-Authenticate"
	HeaderProxyAuthorization              = "Proxy-Authorization"
	HeaderRange                           = "Range"
	HeaderReferer                         = "Referer"
	HeaderRetryAfter                      = "Retry-After"
	HeaderServer                          = "Server"
	HeaderTE                              = "TE"
	HeaderTrailer                         = "Trailer"
	HeaderTransferEncoding                = "Transfer-Encoding"
	HeaderUpgrade                         = "Upgrade"
	HeaderUserAgent                       = "User-Agent"
	HeaderVary                            = "Vary"
	HeaderVia                             = "Via"
	HeaderWarning                         = "Warning"
	HeaderWWWAuthenticate                 = "WWW-Authenticate"
	HeaderXForwardedFor                   = "X-Forwarded-For"
	HeaderXForwardedProto                 = "X-Forwarded-Proto"
	HeaderXForwardedHost                  = "X-Forwarded-Host"
	HeaderXRealIP                         = "X-Real-IP"
	HeaderXRequestID                      = "X-Request-ID"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
)

// Common MIME types
const (
	MimeTypeJSON                  = "application/json"
	MimeTypeXML                   = "application/xml"
	MimeTypeForm                  = "application/x-www-form-urlencoded"
	MimeTypeMultipartForm         = "multipart/form-data"
	MimeTypeOctetStream           = "application/octet-stream"
	MimeTypeTextPlain             = "text/plain"
	MimeTypeTextHTML              = "text/html"
	MimeTypeTextCSS               = "text/css"
	MimeTypeTextJavaScript        = "text/javascript"
	MimeTypeApplicationJavaScript = "application/javascript"
	MimeTypeImageJPEG             = "image/jpeg"
	MimeTypeImagePNG              = "image/png"
	MimeTypeImageGIF              = "image/gif"
	MimeTypeImageSVG              = "image/svg+xml"
	MimeTypeImageWebP             = "image/webp"
	MimeTypeVideoMP4              = "video/mp4"
	MimeTypeVideoWebM             = "video/webm"
	MimeTypeAudioMP3              = "audio/mpeg"
	MimeTypeAudioWAV              = "audio/wav"
	MimeTypeAudioOGG              = "audio/ogg"
	MimeTypeFontWOFF              = "font/woff"
	MimeTypeFontWOFF2             = "font/woff2"
	MimeTypeFontTTF               = "font/ttf"
	MimeTypeFontOTF               = "font/otf"
)

// Default timeout values
const (
	// DefaultRequestTimeout is the default timeout for HTTP requests
	DefaultRequestTimeout = 30 * time.Second

	// DefaultServerReadTimeout is the default read timeout for HTTP server
	DefaultServerReadTimeout = 10 * time.Second

	// DefaultServerWriteTimeout is the default write timeout for HTTP server
	DefaultServerWriteTimeout = 10 * time.Second

	// DefaultServerIdleTimeout is the default idle timeout for HTTP server
	DefaultServerIdleTimeout = 120 * time.Second

	// DefaultKeepAliveTimeout is the default keep-alive timeout
	DefaultKeepAliveTimeout = 75 * time.Second
)

// HTTP constants
const (
	// DefaultHTTPPort is the default HTTP port
	DefaultHTTPPort = 80

	// DefaultHTTPSPort is the default HTTPS port
	DefaultHTTPSPort = 443

	// MaxHeaderSize is the maximum size of HTTP headers
	MaxHeaderSize = 1 << 20 // 1MB

	// MaxRequestBodySize is the maximum size of request body
	MaxRequestBodySize = 10 << 20 // 10MB

	// HTTPSeparator is the HTTP line separator
	HTTPSeparator = "\r\n"

	// HTTPHeaderSeparator separates header name and value
	HTTPHeaderSeparator = ": "

	// HTTPVersionPrefix is the prefix for HTTP version
	HTTPVersionPrefix = "HTTP/"
)

// StatusText returns the status text for the given status code
func StatusText(code StatusCode) string {
	switch code {
	case StatusContinue:
		return "Continue"
	case StatusSwitchingProtocols:
		return "Switching Protocols"
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNonAuthoritativeInfo:
		return "Non-Authoritative Information"
	case StatusNoContent:
		return "No Content"
	case StatusResetContent:
		return "Reset Content"
	case StatusPartialContent:
		return "Partial Content"
	case StatusMultipleChoices:
		return "Multiple Choices"
	case StatusMovedPermanently:
		return "Moved Permanently"
	case StatusFound:
		return "Found"
	case StatusSeeOther:
		return "See Other"
	case StatusNotModified:
		return "Not Modified"
	case StatusUseProxy:
		return "Use Proxy"
	case StatusTemporaryRedirect:
		return "Temporary Redirect"
	case StatusPermanentRedirect:
		return "Permanent Redirect"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusPaymentRequired:
		return "Payment Required"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusNotAcceptable:
		return "Not Acceptable"
	case StatusProxyAuthRequired:
		return "Proxy Authentication Required"
	case StatusRequestTimeout:
		return "Request Timeout"
	case StatusConflict:
		return "Conflict"
	case StatusGone:
		return "Gone"
	case StatusLengthRequired:
		return "Length Required"
	case StatusPreconditionFailed:
		return "Precondition Failed"
	case StatusRequestEntityTooLarge:
		return "Request Entity Too Large"
	case StatusRequestURITooLong:
		return "Request URI Too Long"
	case StatusUnsupportedMediaType:
		return "Unsupported Media Type"
	case StatusRequestedRangeNotSatisfiable:
		return "Requested Range Not Satisfiable"
	case StatusExpectationFailed:
		return "Expectation Failed"
	case StatusTeapot:
		return "I'm a teapot"
	case StatusMisdirectedRequest:
		return "Misdirected Request"
	case StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case StatusLocked:
		return "Locked"
	case StatusFailedDependency:
		return "Failed Dependency"
	case StatusTooEarly:
		return "Too Early"
	case StatusUpgradeRequired:
		return "Upgrade Required"
	case StatusPreconditionRequired:
		return "Precondition Required"
	case StatusTooManyRequests:
		return "Too Many Requests"
	case StatusRequestHeaderFieldsTooLarge:
		return "Request Header Fields Too Large"
	case StatusUnavailableForLegalReasons:
		return "Unavailable For Legal Reasons"
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusGatewayTimeout:
		return "Gateway Timeout"
	case StatusHTTPVersionNotSupported:
		return "HTTP Version Not Supported"
	case StatusVariantAlsoNegotiates:
		return "Variant Also Negotiates"
	case StatusInsufficientStorage:
		return "Insufficient Storage"
	case StatusLoopDetected:
		return "Loop Detected"
	case StatusNotExtended:
		return "Not Extended"
	case StatusNetworkAuthenticationRequired:
		return "Network Authentication Required"
	default:
		return "Unknown Status Code"
	}
}

// IsInformational returns true if the status code is informational (1xx)
func IsInformational(code StatusCode) bool {
	return code >= 100 && code < 200
}

// IsSuccess returns true if the status code indicates success (2xx)
func IsSuccess(code StatusCode) bool {
	return code >= 200 && code < 300
}

// IsRedirection returns true if the status code indicates redirection (3xx)
func IsRedirection(code StatusCode) bool {
	return code >= 300 && code < 400
}

// IsClientError returns true if the status code indicates client error (4xx)
func IsClientError(code StatusCode) bool {
	return code >= 400 && code < 500
}

// IsServerError returns true if the status code indicates server error (5xx)
func IsServerError(code StatusCode) bool {
	return code >= 500 && code < 600
}

// IsError returns true if the status code indicates an error (4xx or 5xx)
func IsError(code StatusCode) bool {
	return IsClientError(code) || IsServerError(code)
}