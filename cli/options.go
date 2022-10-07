package cli

import (
	"better-docker-ps/langext/term"
	"time"
)

type Options struct {
	Version     bool
	Help        bool
	Socket      string
	Quiet       bool
	Verbose     bool
	OutputColor bool
	TimeZone    *time.Location
	Input       *string
	All         bool
	WithSize    bool
	Filter      *map[string]string
}

func DefaultCLIOptions() Options {
	return Options{
		Version:     false,
		Help:        false,
		Quiet:       false,
		Verbose:     false,
		OutputColor: term.SupportsColors(),
		TimeZone:    time.Local,
		Socket:      "/var/run/docker.sock",
		Input:       nil,
		All:         false,
		WithSize:    false,
	}
}
