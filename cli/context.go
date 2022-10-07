package cli

import (
	pserr "better-docker-ps/fferr"
	"better-docker-ps/langext"
	"better-docker-ps/langext/term"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

type PSContext struct {
	context.Context
	Opt Options
}

func (c PSContext) PrintPrimaryOutput(msg string) {
	if c.Opt.Quiet {
		return
	}

	c.printPrimaryRaw(msg + "\n")
}

func (c PSContext) PrintFatalMessage(msg string) {
	if c.Opt.Quiet {
		return
	}

	c.printErrorRaw(msg + "\n")
}

func (c PSContext) PrintFatalError(e error) {
	if c.Opt.Quiet {
		return
	}

	c.printErrorRaw(pserr.FormatError(e, c.Opt.Verbose) + "\n")
}

func (c PSContext) PrintErrorMessage(msg string) {
	if c.Opt.Quiet {
		return
	}

	c.printErrorRaw(msg + "\n")
}

func (c PSContext) PrintVerbose(msg string) {
	if c.Opt.Quiet || !c.Opt.Verbose {
		return
	}

	c.printVerboseRaw(msg + "\n")
}

func (c PSContext) PrintVerboseHeader(msg string) {
	if c.Opt.Quiet || !c.Opt.Verbose {
		return
	}

	c.printVerboseRaw("\n")
	c.printVerboseRaw("========================================" + "\n")
	c.printVerboseRaw(msg + "\n")
	c.printVerboseRaw("========================================" + "\n")
	c.printVerboseRaw("\n")
}

func (c PSContext) PrintVerboseKV(key string, vval any) {
	if c.Opt.Quiet || !c.Opt.Verbose {
		return
	}

	termlen := 236
	keylen := 28

	var val = ""
	switch v := vval.(type) {
	case []byte:
		val = hex.EncodeToString(v)
	case string:
		val = v
	case time.Time:
		val = v.In(c.Opt.TimeZone).Format(time.RFC3339Nano)
	default:
		val = fmt.Sprintf("%v", v)
	}

	if len(val) > (termlen-keylen-4) || strings.Contains(val, "\n") {

		c.printVerboseRaw(key + " :=\n" + val + "\n")

	} else {

		padkey := langext.StrPadRight(key, " ", keylen)
		c.printVerboseRaw(padkey + " := " + val + "\n")

	}
}

func (c PSContext) printPrimaryRaw(msg string) {
	if c.Opt.Quiet {
		return
	}

	writeStdout(msg)
}

func (c PSContext) printErrorRaw(msg string) {
	if c.Opt.Quiet {
		return
	}

	if c.Opt.OutputColor {
		writeStderr(term.Red(msg))
	} else {
		writeStderr(msg)
	}
}

func (c PSContext) printVerboseRaw(msg string) {
	if c.Opt.Quiet {
		return
	}

	if c.Opt.OutputColor {
		writeStdout(term.Gray(msg))
	} else {
		writeStdout(msg)
	}
}

func writeStdout(msg string) {
	_, err := os.Stdout.WriteString(msg)
	if err != nil {
		panic("failed to write to stdout: " + err.Error())
	}
}

func writeStderr(msg string) {
	_, err := os.Stderr.WriteString(msg)
	if err != nil {
		panic("failed to write to stdout: " + err.Error())
	}
}

func NewContext(opt Options) (*PSContext, error) {
	return &PSContext{
		Context: context.Background(),
		Opt:     opt,
	}, nil
}

func NewEarlyContext() *PSContext {
	return &PSContext{
		Context: context.Background(),
		Opt:     DefaultCLIOptions(),
	}
}

func (c PSContext) Finish() {
	// ...
}
