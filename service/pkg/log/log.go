// Package log provides functions for configuring global logging.
package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureGlobalLogging parses a given log level and sets it globally.
func ConfigureGlobalLogging(level, service, version string) {
	// parse and set the global log level
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	// NOTE: global logger settings can be found here: https://github.com/rs/zerolog#global-settings
	zerolog.SetGlobalLevel(lvl)

	// configure global logger defaults
	log.Logger = log.Logger.With().Fields(map[string]interface{}{
		"service": service,
		"version": version,
	}).Logger()
}
