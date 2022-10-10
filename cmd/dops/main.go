package main

import (
	"better-docker-ps/cli"
	"better-docker-ps/consts"
	pserr "better-docker-ps/fferr"
	"better-docker-ps/impl"
	"fmt"
	"os"
	"runtime/debug"
)

// Inspiration: https://github.com/moby/moby/issues/7477

func main() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%v\n\n%s", err, string(debug.Stack())))
			os.Exit(consts.ExitcodePanic)
		}
	}()

	opt, err := cli.ParseCommandline()
	if err != nil {
		ctx := cli.NewEarlyContext()
		ctx.PrintFatalError(err)
		os.Exit(pserr.GetExitCode(err, consts.ExitcodeCLIParse))
		return
	}

	ctx, err := cli.NewContext(opt)
	if err != nil {
		ctx.PrintFatalError(err)
		os.Exit(pserr.GetExitCode(err, consts.ExitcodeError))
		return
	}

	defer ctx.Finish()

	if opt.Version {
		ctx.PrintPrimaryOutput("better-docker-ps v" + consts.BETTER_DOCKER_PS_VERSION)
		os.Exit(0)
		return
	}

	if opt.Help {
		printHelp(ctx)
		os.Exit(0)
		return
	}

	err = impl.Execute(ctx)
	if err != nil {
		ctx.PrintFatalError(err)
		os.Exit(pserr.GetExitCode(err, consts.ExitcodeError))
		return
	}

	os.Exit(0)
}

func printHelp(ctx *cli.PSContext) {
	ctx.PrintPrimaryOutput("better-docker-ps")
	ctx.PrintPrimaryOutput("")
	ctx.PrintPrimaryOutput("Usage:")
	ctx.PrintPrimaryOutput("  dops [OPTIONS]                     List docker container")
	ctx.PrintPrimaryOutput("")
	ctx.PrintPrimaryOutput("Options (default):")
	ctx.PrintPrimaryOutput("  -h, --help                         Show this screen.")
	ctx.PrintPrimaryOutput("  --version                          Show version.")
	ctx.PrintPrimaryOutput("  --all , -a                         Show all containers (default shows just running)")
	ctx.PrintPrimaryOutput("  --filter <ftr>, -f <ftr>           Filter output based on conditions provided")
	ctx.PrintPrimaryOutput("  --format <fmt>                     Pretty-print containers using a Go template")
	ctx.PrintPrimaryOutput("  --last , -n                        Show n last created containers (includes all states)")
	ctx.PrintPrimaryOutput("  --latest , -l                      Show the latest created container (includes all states)")
	ctx.PrintPrimaryOutput("  --no-trunc                         Don't truncate output")
	ctx.PrintPrimaryOutput("  --quiet , -q                       Only display container IDs")
	ctx.PrintPrimaryOutput("  --size , -s                        Display total file sizes")
	ctx.PrintPrimaryOutput("")
	ctx.PrintPrimaryOutput("Options (extra | do not exist in `docker ps`):")
	ctx.PrintPrimaryOutput("  --silent                           Do not print any output")
	ctx.PrintPrimaryOutput("  --timezone                         Specify the timezone for date outputs")
	ctx.PrintPrimaryOutput("  --color <true|false>               Enable/Disable terminal color output")
	ctx.PrintPrimaryOutput("  --no-color                         Disable terminal color output")
	ctx.PrintPrimaryOutput("  --socket <filepath>                Specify the docker socket location (Default: /var/run/docker.sock)")
	ctx.PrintPrimaryOutput("  --timeformat <go-time-fmt>         Specify the datetime output format (golang syntax)")
	ctx.PrintPrimaryOutput("  --no-header                        Do not print the table header")
	ctx.PrintPrimaryOutput("")
	ctx.PrintPrimaryOutput("Available --format keys (default):")
	ctx.PrintPrimaryOutput("  {{.ID}                             Container ID")
	ctx.PrintPrimaryOutput("  {{.Image}                          Image ID")
	ctx.PrintPrimaryOutput("  {{.Command}                        Quoted command")
	ctx.PrintPrimaryOutput("  {{.CreatedAt}                      Time when the container was created.")
	ctx.PrintPrimaryOutput("  {{.RunningFor}                     Elapsed time since the container was started.")
	ctx.PrintPrimaryOutput("  {{.Ports}                          Exposed ports.")
	ctx.PrintPrimaryOutput("  {{.State}                          Container status")
	ctx.PrintPrimaryOutput("  {{.Status}                         Container status with details")
	ctx.PrintPrimaryOutput("  {{.Size}                           Container disk size.")
	ctx.PrintPrimaryOutput("  {{.Names}                          Container names.")
	ctx.PrintPrimaryOutput("  {{.Labels}                         All labels assigned to the container.")
	ctx.PrintPrimaryOutput("  {{.Label}                          [!] Unsupported")
	ctx.PrintPrimaryOutput("  {{.Mounts}                         Names of the volumes mounted in this container.")
	ctx.PrintPrimaryOutput("  {{.Networks}                       Names of the networks attached to this container.")
	ctx.PrintPrimaryOutput("")
	ctx.PrintPrimaryOutput("Available --format keys (extra | do not exist in `docker ps`):")
	ctx.PrintPrimaryOutput("  {{.ImageName}                      Image ID (without tag and registry)")
	ctx.PrintPrimaryOutput("  {{.ImageTag}, {{.Tag}              Image Tag")
	ctx.PrintPrimaryOutput("  {{.ImageRegistry}, {{.Registry}    Image Registry")
	ctx.PrintPrimaryOutput("  {{.ShortCommand}                   Command without arguments")
	ctx.PrintPrimaryOutput("  {{.LabelsKeys}                     All labels assigned to the container (keys only)")
	ctx.PrintPrimaryOutput("  {{.IP}                             Internal IP Address")
	ctx.PrintPrimaryOutput("")
}
