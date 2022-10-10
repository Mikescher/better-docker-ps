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
	TimeFormat  string
	Input       *string
	All         bool
	WithSize    bool
	Filter      *map[string]string
	OnlyIDs     bool
	Limit       int
	Format      string
	PrintHeader bool
	Truncate    bool
}

func DefaultCLIOptions() Options {
	return Options{
		Version:     false,
		Help:        false,
		Quiet:       false,
		Verbose:     false,
		OutputColor: term.SupportsColors(),
		TimeZone:    time.Local,
		TimeFormat:  "2006-01-02 15:04:05 Z07:00 MST",
		Socket:      "/var/run/docker.sock",
		Input:       nil,
		All:         false,
		WithSize:    false,
		OnlyIDs:     false,
		Limit:       -1,
		Format:      "table {{.ID}}\t{{.Names}}\t{{.Tag}}\t{{.CreatedAt}}\t{{.State}}\t{{.Status}}\t{{.Ports}}\t{{.IP}}",
		PrintHeader: true,
		Truncate:    true,
	}
}
