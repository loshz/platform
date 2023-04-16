package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc/credentials"
)

// LoadGRPCCreds takes the paths of a server cert/key and client ca and returns
// TLS transport credentials. TLS 1.3 is the minimum version.
func LoadGRPCCreds(crt, key, ca string) (credentials.TransportCredentials, error) {
	// Load server key pair.
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}

	// Load client CA.
	data, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		return nil, err
	}

	// Configure TLS options.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    capool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
}
