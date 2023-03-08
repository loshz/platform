package main

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/loshz/platform/pkg/service"
)

func main() {
	s := service.New("trafficd")
	s.Run(run)
}

func run(s *service.Service) error {
	t := time.NewTicker(10 * time.Second)
	go func() {
		for range t.C {
			log.Info().Msg("ticking...")
		}
	}()

	return nil
}
