package service

import (
	"fmt"

	"github.com/loshz/platform/internal/credentials"
)

// LoadCredentials attempts to load service credentials from config values into
// the credentials store.
// As this method is intended to be ran before service startup, errors are treated
// as fatal.
func (s *Service) LoadCredentials(creds ...credentials.Credential) {
	for _, cred := range creds {
		var err error

		switch cred {
		case credentials.GrpcClient:
			s.LoadGrpcClientConfig()
			err = s.Creds().LoadGrpcClientCreds(s.Config())
		case credentials.GrpcServer:
			s.LoadGrpcServerConfig()
			err = s.Creds().LoadGrpcServerCreds(s.Config())
		}

		if err != nil {
			panic(fmt.Errorf("error loading service credentials: %w", err))
		}
	}
}
