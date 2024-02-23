package config

const (
	// Service config.
	KeyServiceLogLevel         = "service.log.level"
	KeyServiceShutdownTimeout  = "service.shutdown.timeout"
	KeyServiceDiscoveryEnabled = "service.discovery.enabled"
	KeyServiceDiscoveryAddr    = "service.discovery.addr"
	KeyServiceRegisterInt      = "service.register.interval"

	// HTTPS/S server config.
	KeyHttpPort         = "http.port"
	KeyHttpReadTimeout  = "http.read.timeout"
	KeyHttpWriteTimeout = "http.write.timeout"
	KeyHttpIdleTimeout  = "http.idle.timeout"

	// gRPC TLS config.
	KeyGrpcTlsCA = "grpc.tls.ca"

	// gRPC server config.
	KeyGrpcServerPort        = "grpc.server.port"
	KeyGrpcServerCert        = "grpc.server.cert"
	KeyGrpcServerKey         = "grpc.server.key"
	KeyGrpcServerConnTimeout = "grpc.server.conn.timeout"

	// gRPC client config.
	KeyGrpcClientCert    = "grpc.client.cert"
	KeyGrpcClientKey     = "grpc.client.key"
	KeyGrpcClientTimeout = "grpc.client.timeout"
)
