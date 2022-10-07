package main

import (
	"better-docker-ps/cli"
	"better-docker-ps/cli/parser"
	"better-docker-ps/consts"
	pserr "better-docker-ps/fferr"
	"better-docker-ps/impl"
	"fmt"
	"os"
	"runtime/debug"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%v\n\n%s", err, string(debug.Stack())))
			os.Exit(consts.ExitcodePanic)
		}
	}()

	opt, err := parser.ParseCommandline()
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

	err = impl.Execute(ctx)
	if err != nil {
		ctx.PrintFatalError(err)
		os.Exit(pserr.GetExitCode(err, consts.ExitcodeError))
		return
	}

	os.Exit(0)
}
