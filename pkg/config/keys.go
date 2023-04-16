package config

const (
	// The level at which the service should write application logs.
	KeyLogLevel = "log.level"

	// Service related config.
	KeyServiceStartupTimeout = "service.startup.timeout"

	// HTTPS/S server related fields/timeouts.
	KeyHTTPPort         = "http.port"
	KeyHTTPReadTimeout  = "http.read.timeout"
	KeyHTTPWriteTimeout = "http.write.timeout"
	KeyHTTPIdleTimeout  = "http.idle.timeout"

	// gRPC server related fields/timeouts.
	KeyGRPCServerPort        = "grpc.server.port"
	KeyGRPCServerCert        = "grpc.server.cert"
	KeyGRPCServerKey         = "grpc.server.key"
	KeyGRPCServerConnTimeout = "grpc.server.conn.timeout"
)
