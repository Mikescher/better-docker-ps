package parser

import (
	"better-docker-ps/cli"
	"better-docker-ps/fferr"
	"better-docker-ps/langext"
	"github.com/joomcode/errorx"
	"os"
	"strings"
	"time"
)

func ParseCommandline() (cli.Options, error) {
	o, err := parseCommandlineInternal()
	if err != nil {
		return cli.Options{}, errorx.Decorate(err, "failed to parse commandline")
	}
	return o, nil
}

func parseCommandlineInternal() (cli.Options, error) {
	unprocessedArgs := os.Args[1:]

	allOptionArguments := make([]cli.ArgumentTuple, 0)

	for len(unprocessedArgs) > 0 {
		arg := unprocessedArgs[0]
		unprocessedArgs = unprocessedArgs[1:]

		if strings.HasPrefix(arg, "--") {

			arg = arg[2:]

			if strings.Contains(arg, "=") {
				key := arg[0:strings.Index(arg, "=")]
				val := arg[strings.Index(arg, "=")+1:]

				if len(key) <= 1 {
					return cli.Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
				}

				allOptionArguments = append(allOptionArguments, cli.ArgumentTuple{Key: key, Value: langext.Ptr(val)})
				continue
			} else {

				key := arg

				if len(key) <= 1 {
					return cli.Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
				}

				if len(unprocessedArgs) == 0 || strings.HasPrefix(unprocessedArgs[0], "-") {
					allOptionArguments = append(allOptionArguments, cli.ArgumentTuple{Key: key, Value: nil})
					continue
				} else {
					val := unprocessedArgs[0]
					unprocessedArgs = unprocessedArgs[1:]
					allOptionArguments = append(allOptionArguments, cli.ArgumentTuple{Key: key, Value: langext.Ptr(val)})
					continue
				}

			}

		} else if strings.HasPrefix(arg, "-") {

			arg = arg[1:]

			if len(arg) > 1 {
				for i := 1; i < len(arg); i++ {
					allOptionArguments = append(allOptionArguments, cli.ArgumentTuple{Key: arg[i : i+1], Value: nil})
				}
				continue
			}

			key := arg

			if key == "" {
				return cli.Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
			}

			if len(unprocessedArgs) == 0 || strings.HasPrefix(unprocessedArgs[0], "-") {
				allOptionArguments = append(allOptionArguments, cli.ArgumentTuple{Key: key, Value: nil})
				continue
			} else {
				val := unprocessedArgs[0]
				unprocessedArgs = unprocessedArgs[1:]
				allOptionArguments = append(allOptionArguments, cli.ArgumentTuple{Key: key, Value: langext.Ptr(val)})
				continue
			}

		} else {
			return cli.Options{}, pserr.DirectOutput.New("Unknown/Misplaced argument: " + arg)
		}
	}

	// Process common options

	opt := cli.DefaultCLIOptions()

	for _, arg := range allOptionArguments {

		if (arg.Key == "h" || arg.Key == "help") && arg.Value == nil {
			return cli.Options{Help: true}, nil
		}

		if arg.Key == "version" && arg.Value == nil {
			return cli.Options{Version: true}, nil
		}

		if arg.Key == "version" && arg.Value == nil {
			return cli.Options{Version: true}, nil
		}

		if (arg.Key == "v" || arg.Key == "verbose") && arg.Value == nil {
			opt.Verbose = true
			continue
		}

		if (arg.Key == "q" || arg.Key == "quiet") && arg.Value == nil {
			opt.Quiet = true
			continue
		}

		if arg.Key == "timezone" && arg.Value != nil {
			loc, err := time.LoadLocation(*arg.Value)
			if err != nil {
				return cli.Options{}, pserr.DirectOutput.New("Unknown timezone: " + *arg.Value)
			}
			opt.TimeZone = loc
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
				return cli.Options{}, pserr.DirectOutput.New("Filter argument must have a key and a value (a=b): " + arg.Key)
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

		return cli.Options{}, pserr.DirectOutput.New("Unknown argument: " + arg.Key)
	}

	return opt, nil
}
