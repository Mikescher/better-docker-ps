package cli

import (
	"gogs.mikescher.com/BlackForestBytes/goext/termext"
	"time"
)

type SortDirection string

const (
	SortASC  SortDirection = "ASC"
	SortDESC SortDirection = "DESC"
)

type Options struct {
	Version          bool
	Help             bool
	Socket           string
	Quiet            bool
	Verbose          bool
	OutputColor      bool
	TimeZone         *time.Location
	TimeFormat       string
	TimeFormatHeader string
	Input            *string
	All              bool
	WithSize         bool
	Filter           *map[string]string
	Limit            int
	DefaultFormat    bool
	Format           []string // if more than 1 value, we use the later values as fallback for too-small terminal
	PrintHeader      bool
	PrintHeaderLines bool
	Truncate         bool
	SortColumns      []string
	SortDirection    []SortDirection
	WatchInterval    *time.Duration
}

func DefaultCLIOptions() Options {
	return Options{
		Version:          false,
		Help:             false,
		Quiet:            false,
		Verbose:          false,
		OutputColor:      termext.SupportsColors(),
		TimeZone:         time.Local,
		TimeFormatHeader: "Z07:00 MST",
		TimeFormat:       "2006-01-02 15:04:05",
		Socket:           "/var/run/docker.sock",
		Input:            nil,
		All:              false,
		WithSize:         false,
		Limit:            -1,
		DefaultFormat:    true,
		Format: []string{
			"table {{.ID}}\\t{{.Names}}\\t{{.ImageName}}\\t{{.Tag}}\\t{{.ShortCommand}}\\t{{.CreatedAt}}\\t{{.State}}\\t{{.Status}}\\t{{.LongPublishedPorts}}\\t{{.Networks}}\\t{{.IP}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.ImageName}}\\t{{.Tag}}\\t{{.ShortCommand}}\\t{{.CreatedAt}}\\t{{.State}}\\t{{.Status}}\\t{{.LongPublishedPorts}}\\t{{.IP}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.ImageName}}\\t{{.Tag}}\\t{{.CreatedAt}}\\t{{.State}}\\t{{.Status}}\\t{{.LongPublishedPorts}}\\t{{.IP}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.ImageName}}\\t{{.Tag}}\\t{{.CreatedAt}}\\t{{.State}}\\t{{.Status}}\\t{{.PublishedPorts}}\\t{{.IP}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.ImageName}}\\t{{.Tag}}\\t{{.CreatedAt}}\\t{{.State}}\\t{{.Status}}\\t{{.PublishedPorts}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.ImageName}}\\t{{.Tag}}\\t{{.State}}\\t{{.Status}}\\t{{.PublishedPorts}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.Tag}}\\t{{.State}}\\t{{.Status}}\\t{{.PublishedPorts}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.Tag}}\\t{{.State}}\\t{{.Status}}\\t{{.ShortPublishedPorts}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.Tag}}\\t{{.State}}\\t{{.Status}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.State}}\\t{{.Status}}",
			"table {{.ID}}\\t{{.Names}}\\t{{.State}}",
			"table {{.Names}}\\t{{.State}}",
			"table {{.Names}}",
			"table {{.ID}}",
		},
		PrintHeader:      true,
		PrintHeaderLines: true,
		Truncate:         true,
		SortColumns:      make([]string, 0),
		SortDirection:    make([]SortDirection, 0),
		WatchInterval:    nil,
	}
}
