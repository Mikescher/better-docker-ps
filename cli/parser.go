package cli

import (
	"better-docker-ps/pserr"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joomcode/errorx"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
)

func ParseCommandline() (Options, error) {
	o, err := parseCommandlineInternal()
	if err != nil {
		return Options{}, errorx.Decorate(err, "failed to parse commandline")
	}
	return o, nil
}

func parseCommandlineInternal() (Options, error) {
	unprocessedArgs := os.Args[1:]

	allOptionArguments := make([]ArgumentTuple, 0)

	for len(unprocessedArgs) > 0 {
		arg := unprocessedArgs[0]
		unprocessedArgs = unprocessedArgs[1:]

		if strings.HasPrefix(arg, "--") {

			arg = arg[2:]

			if strings.Contains(arg, "=") {
				key := arg[0:strings.Index(arg, "=")]
				val := arg[strings.Index(arg, "=")+1:]

				if len(key) <= 1 {
					return Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
				}

				allOptionArguments = append(allOptionArguments, ArgumentTuple{Key: key, Value: langext.Ptr(val)})
				continue
			} else {

				key := arg

				if len(key) <= 1 {
					return Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
				}

				if len(unprocessedArgs) == 0 || strings.HasPrefix(unprocessedArgs[0], "-") {
					allOptionArguments = append(allOptionArguments, ArgumentTuple{Key: key, Value: nil})
					continue
				} else {
					val := unprocessedArgs[0]
					unprocessedArgs = unprocessedArgs[1:]
					allOptionArguments = append(allOptionArguments, ArgumentTuple{Key: key, Value: langext.Ptr(val)})
					continue
				}

			}

		} else if strings.HasPrefix(arg, "-") {

			arg = arg[1:]

			if len(arg) > 1 {
				for i := 0; i < len(arg); i++ {
					allOptionArguments = append(allOptionArguments, ArgumentTuple{Key: arg[i : i+1], Value: nil})
				}
				continue
			}

			key := arg

			if key == "" {
				return Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
			}

			if len(unprocessedArgs) == 0 || strings.HasPrefix(unprocessedArgs[0], "-") {
				allOptionArguments = append(allOptionArguments, ArgumentTuple{Key: key, Value: nil})
				continue
			} else {
				val := unprocessedArgs[0]
				unprocessedArgs = unprocessedArgs[1:]
				allOptionArguments = append(allOptionArguments, ArgumentTuple{Key: key, Value: langext.Ptr(val)})
				continue
			}

		} else {
			return Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
		}
	}

	// Process common options

	opt := DefaultCLIOptions()

	for _, arg := range allOptionArguments {

		if (arg.Key == "h" || arg.Key == "help") && arg.Value == nil {
			return Options{Help: true}, nil
		}

		if arg.Key == "version" && arg.Value == nil {
			return Options{Version: true}, nil
		}

		if (arg.Key == "v" || arg.Key == "verbose") && arg.Value == nil {
			opt.Verbose = true
			continue
		}

		if (arg.Key == "silent") && arg.Value == nil {
			opt.Quiet = true
			continue
		}

		if (arg.Key == "q" || arg.Key == "quiet") && arg.Value == nil {
			opt.Format = []string{"idlist"}
			opt.DefaultFormat = false
			continue
		}

		if arg.Key == "timezone" && arg.Value != nil {
			loc, err := time.LoadLocation(*arg.Value)
			if err != nil {
				return Options{}, pserr.DirectOutput.New("Unknown timezone: " + *arg.Value)
			}
			opt.TimeZone = loc
			continue
		}

		if arg.Key == "timeformat" && arg.Value != nil {
			opt.TimeFormat = *arg.Value
			opt.TimeFormatHeader = ""
			continue
		}

		if arg.Key == "timeformat-header" && arg.Value != nil {
			opt.TimeFormatHeader = *arg.Value
			continue
		}

		if arg.Key == "color" && arg.Value != nil && (strings.ToLower(*arg.Value) == "true" || *arg.Value == "1") {
			opt.OutputColor = true
			continue
		}

		if arg.Key == "color" && arg.Value != nil && (strings.ToLower(*arg.Value) == "false" || *arg.Value == "0") {
			opt.OutputColor = false
			continue
		}

		if arg.Key == "no-color" && arg.Value == nil {
			opt.OutputColor = false
			continue
		}

		if (arg.Key == "socket") && arg.Value != nil {
			opt.Socket = *arg.Value
			continue
		}

		if (arg.Key == "input") && arg.Value != nil {
			// used for testing
			opt.Input = langext.Ptr(*arg.Value)
			continue
		}

		if (arg.Key == "all" || arg.Key == "a") && arg.Value == nil {
			opt.All = true
			continue
		}

		if (arg.Key == "size") && arg.Value == nil {
			opt.WithSize = true
			continue
		}

		if (arg.Key == "filter") && arg.Value != nil {
			spl := strings.SplitN(*arg.Value, "=", 2)
			if len(spl) != 2 {
				return Options{}, pserr.DirectOutput.New("Filter argument must have a key and a value (a=b): " + arg.Key)
			}
			if opt.Filter == nil {
				_v := make(map[string]string)
				opt.Filter = &_v
			}
			filter := *opt.Filter
			filter[spl[0]] = spl[1]
			opt.Filter = &filter
			continue
		}

		if (arg.Key == "format") && arg.Value != nil {
			if opt.DefaultFormat {
				opt.Format = make([]string, 0)
			}
			opt.Format = append(opt.Format, *arg.Value)
			opt.DefaultFormat = false
			continue
		}

		if (arg.Key == "last" || arg.Key == "n") && arg.Value != nil {
			if v, err := strconv.ParseInt(*arg.Value, 10, 32); err == nil {
				opt.Limit = int(v)
				opt.All = true
				continue
			}
			return Options{}, pserr.DirectOutput.New("Failed to parse number argument '--last': '" + *arg.Value + "'")
		}

		if (arg.Key == "latest" || arg.Key == "l") && arg.Value != nil {
			opt.Limit = 1
			opt.All = true
			continue
		}

		if (arg.Key == "no-trunc" || arg.Key == "no-truncate") && arg.Value == nil {
			opt.Truncate = false
			continue
		}

		if (arg.Key == "no-header") && arg.Value == nil {
			opt.PrintHeader = false
			continue
		}

		if (arg.Key == "simple-header") && arg.Value == nil {
			opt.PrintHeaderLines = false
			continue
		}

		return Options{}, pserr.DirectOutput.New("Unknown argument: " + arg.Key)
	}

	return opt, nil
}
