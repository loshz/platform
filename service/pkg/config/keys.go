package config

const (
	// The level at which the service should write application
	KeyLogLevel = "log.level"

	// The port to expose the HTTP/S server on
	KeyHTTPPort = "http.port"

	// HTTP/S server timeouts
	KeyHTTPReadTimeout  = "http.read.timeout"
	KeyHTTPWriteTimeout = "http.write.timeout"
	KeyHTTPIdleTimeout  = "http.idle.timeout"
)
