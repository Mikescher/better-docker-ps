package cli

import (
	"encoding/json"
	"regexp"
	"time"

	"git.blackforestbytes.com/BlackForestBytes/goext/cmdext"
	"git.blackforestbytes.com/BlackForestBytes/goext/termext"
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
		Socket:           "auto",
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

func (o Options) GetSocket() string {
	const defaultSocket = "/var/run/docker.sock"

	if o.Socket != "auto" {
		return o.Socket
	}

	res, err := cmdext.Runner("docker").Arg("context").Arg("list").Arg("--format").Arg("json").Timeout(10 * time.Second).FailOnTimeout().FailOnExitCode().Run()
	if err != nil {
		// on error we just return the default socket
		return defaultSocket
	}

	var context dockerContext
	err = json.Unmarshal([]byte(res.StdOut), &context)
	if err != nil {
		// on error we just return the default socket
		return defaultSocket
	}
	if context.Current {
		return context.socket()
	}

	// if we don't have a current context, we just return the default socket
	return defaultSocket
}

type dockerContext struct {
	Name           string
	Description    string
	DockerEndpoint string
	Current        bool
	Error          string
	ContextType    string
}

var unixSocketPrefixPat = regexp.MustCompile("^unix://")

// Get the socket from the docker context line.
//
// This just strips the `unix://` prefix from it if it is there.
func (ctx dockerContext) socket() string {
	return unixSocketPrefixPat.ReplaceAllString(ctx.DockerEndpoint, "")
}
