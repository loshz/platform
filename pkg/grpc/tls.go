package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

// NewTLSConfig takes the paths of a server cert/key and returns
// TLS transport credentials. TLS 1.3 is the minimum version.
func NewTLSConfig(crt, key string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}

	return tlsConfig, nil
}

// NewCertPool takes the path of a CA cert and attempts to parse
// the certificate, adding it to a CA pool if successful.
func NewCertPool(ca string) (*x509.CertPool, error) {
	data, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	capool := x509.NewCertPool()
	if ok := capool.AppendCertsFromPEM(data); !ok {
		return nil, err
	}

	return capool, nil
}

// NewServerTransportCreds creates gRPC Transport Credentials with a client CA
// configured.
func NewServerTransportCreds(crt, key, ca string) (credentials.TransportCredentials, error) {
	capool, err := NewCertPool(ca)
	if err != nil {
		return nil, fmt.Errorf("error loading ca cert pool : %w", err)
	}

	tlsConfig, err := NewTLSConfig(crt, key)
	if err != nil {
		return nil, fmt.Errorf("error loading tls config: %w", err)
	}

	// Set the client CA.
	tlsConfig.ClientCAs = capool

	return credentials.NewTLS(tlsConfig), nil
}

// NewClientTransportCreds creates gRPC Transport Credentials with a root CA
// configured.
func NewClientTransportCreds(crt, key, ca string) (credentials.TransportCredentials, error) {
	capool, err := NewCertPool(ca)
	if err != nil {
		return nil, fmt.Errorf("error loading ca cert pool : %w", err)
	}

	tlsConfig, err := NewTLSConfig(crt, key)
	if err != nil {
		return nil, fmt.Errorf("error loading tls config: %w", err)
	}

	// Set the root CA.
	tlsConfig.RootCAs = capool

	return credentials.NewTLS(tlsConfig), nil
}
