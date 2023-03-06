// Package version exposes public constants used to configure service
// data at build time.
package version

// Version represents the current build number used in `service.Server`.
//
// Use `--ldflags` at build time to set this value.
var Version = "dev"
