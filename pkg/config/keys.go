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

	// gRPC TLS config.
	KeyGRPCTLSCA = "grpc.tls.ca"

	// gRPC server config.
	KeyGRPCServerPort        = "grpc.server.port"
	KeyGRPCServerCert        = "grpc.server.cert"
	KeyGRPCServerKey         = "grpc.server.key"
	KeyGRPCServerConnTimeout = "grpc.server.conn.timeout"

	// gRPC client config.
	KeyGRPCClientCert = "grpc.client.cert"
	KeyGRPCClientKey  = "grpc.client.key"
)
