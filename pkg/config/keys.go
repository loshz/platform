package config

const (
	// Service config.
	KeyLogLevel              = "log.level"
	KeyServiceStartupTimeout = "service.startup.timeout"

	// HTTPS/S server config.
	KeyHTTPPort         = "http.port"
	KeyHTTPReadTimeout  = "http.read.timeout"
	KeyHTTPWriteTimeout = "http.write.timeout"
	KeyHTTPIdleTimeout  = "http.idle.timeout"

	// gRPC server config.
	KeyGRPCServerPort        = "grpc.server.port"
	KeyGRPCServerCert        = "grpc.server.cert"
	KeyGRPCServerKey         = "grpc.server.key"
	KeyGRPCServerConnTimeout = "grpc.server.conn.timeout"

	// gRPC client config.
	KeyGRPCClientCA = "grpc.client.ca"
)
