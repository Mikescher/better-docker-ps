package cli

import (
	"better-docker-ps/pserr"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/termext"
)

type PSContext struct {
	context.Context
	Opt   Options
	Cache map[string]any
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

func (c PSContext) ClearTerminal() {
	fmt.Print("\033[H\033[2J")
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
		writeStderr(termext.Red(msg))
	} else {
		writeStderr(msg)
	}
}

func (c PSContext) printVerboseRaw(msg string) {
	if c.Opt.Quiet {
		return
	}

	if c.Opt.OutputColor {
		writeStdout(termext.Gray(msg))
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
		Cache:   make(map[string]any),
	}, nil
}

func NewEarlyContext() *PSContext {
	return &PSContext{
		Context: context.Background(),
		Opt:     DefaultCLIOptions(),
		Cache:   make(map[string]any),
	}
}

func (c PSContext) Finish() {
	// ...
}

func (c *PSContext) GetIntFromCache(key string, calc func() int) int {
	if v1, ok := c.Cache[key]; ok {
		if v2, ok := v1.(int); ok {
			return v2
		}
		panic(fmt.Sprintf("Wrong type in cache type(%s) = %T  (expected: int)", key, v1))
	}

	val := calc()
	c.Cache[key] = val
	return val
}
