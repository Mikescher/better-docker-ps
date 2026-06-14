package impl

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	"better-docker-ps/printer"
	"better-docker-ps/pserr"
	"encoding/json"
	"fmt"
	"git.blackforestbytes.com/BlackForestBytes/goext/langext"
	"git.blackforestbytes.com/BlackForestBytes/goext/mathext"
	"git.blackforestbytes.com/BlackForestBytes/goext/syncext"
	"golang.org/x/term"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func Execute(ctx *cli.PSContext) error {
	return executeSingle(ctx, false)
}

func Watch(ctx *cli.PSContext, d time.Duration) error {

	sigTermChannel := make(chan os.Signal, 8)
	signal.Notify(sigTermChannel, os.Interrupt, syscall.SIGTERM)

	for {

		err := executeSingle(ctx, true)
		if err != nil {
			return err
		}

		_, isSig := syncext.ReadChannelWithTimeout(sigTermChannel, d)
		if isSig {
			ctx.PrintPrimaryOutput("")
			ctx.PrintPrimaryOutput("Watch canceled with Ctrl+C")
			return nil
		}

	}
}

func executeSingle(ctx *cli.PSContext, clear bool) error {
	for _, f := range ctx.Opt.Format {
		if strings.Contains(f, "{{.Size}}") {
			ctx.Opt.WithSize = true
		}
	}

	jsonraw, err := docker.ListContainer(ctx)
	if err != nil {
		return err
	}

	ctx.PrintVerboseKV("API response", langext.TryPrettyPrintJson(string(jsonraw)))

	var data []docker.ContainerSchema
	err = json.Unmarshal(jsonraw, &data)
	if err != nil {
		return pserr.DirectOutput.Wrap(err, "Failed to decode Docker API response")
	}

	enrichWithInspectData(ctx, data)

	if len(ctx.Opt.SortColumns) > 0 {
		data = doSort(ctx, data, ctx.Opt.SortColumns, ctx.Opt.SortDirection)
	}

	if ctx.Opt.Search != nil {
		data = doSearch(ctx, data, *ctx.Opt.Search)
	}

	for i, v := range ctx.Opt.Format {

		if clear {
			ctx.ClearTerminal()
		}

		ok, err := doOutput(ctx, data, v, i == len(ctx.Opt.Format)-1)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

	}

	return pserr.DirectOutput.New("Missing format specification for output")
}

// enrichWithInspectData fetches additional per-container data from the container-inspect endpoint,
// but only when a requested column actually needs it (currently only the User column).
// The container-list endpoint does not return Config.User, so it has to be queried separately.
func enrichWithInspectData(ctx *cli.PSContext, data []docker.ContainerSchema) {
	if !needsInspectData(ctx) {
		return
	}

	// inspect is not available when the data was loaded from an --input file
	if ctx.Opt.Input != nil {
		return
	}

	for i := range data {
		insp, err := docker.InspectContainer(ctx, data[i].ID)
		if err != nil {
			ctx.PrintVerbose(fmt.Sprintf("Failed to inspect container '%s': %v", data[i].ID, err.Error()))
			continue
		}
		data[i].Config = &insp.Config
	}
}

func needsInspectData(ctx *cli.PSContext) bool {
	for _, f := range ctx.Opt.Format {
		if strings.Contains(f, ".User") || strings.Contains(f, ".Config") {
			return true
		}
	}
	for _, s := range ctx.Opt.SortColumns {
		if s == "User" {
			return true
		}
	}
	return false
}

func doSearch(ctx *cli.PSContext, data []docker.ContainerSchema, needle string) []docker.ContainerSchema {
	needle = strings.ToLower(needle)

	haystackFormat := ""
	for _, f := range ctx.Opt.Format {
		if strings.HasPrefix(f, "table ") {
			haystackFormat = f
			break
		}
	}

	result := make([]docker.ContainerSchema, 0, len(data))
	for _, cont := range data {
		hay := cont.ID + " " + strings.Join(cont.Names, " ") + " " + cont.Image + " " + cont.Command
		if haystackFormat != "" {
			for _, fn := range parseTableDef(haystackFormat) {
				hay += " " + strings.Join(fn(ctx, data, &cont), " ")
			}
		} else if len(ctx.Opt.Format) > 0 {
			hay += " " + replaceSingleLineColumnData(ctx, data, cont, ctx.Opt.Format[0])
		}
		if strings.Contains(strings.ToLower(hay), needle) {
			result = append(result, cont)
		}
	}
	return result
}

func doSort(ctx *cli.PSContext, data []docker.ContainerSchema, skeys []string, sdirs []cli.SortDirection) []docker.ContainerSchema {

	langext.SortSliceStable(data, func(v1, v2 docker.ContainerSchema) bool {

		// return true if v1 < v2

		for i := 0; i < len(skeys); i++ {

			sfn, ok := getSortFun(skeys[i])
			if !ok {
				continue
			}

			cmp := sfn(ctx, &v1, &v2)
			if sdirs[i] == "DESC" {
				cmp = cmp * -1
			}

			if cmp < 0 {
				return true
			} else if cmp > 0 {
				return false
			}
		}

		return false // equals
	})

	return data
}

func doOutput(ctx *cli.PSContext, data []docker.ContainerSchema, format string, force bool) (bool, error) {
	if format == "idlist" {

		for _, v := range data {
			if ctx.Opt.Truncate {
				ctx.PrintPrimaryOutput(v.ID[0:12])
			} else {
				ctx.PrintPrimaryOutput(v.ID)
			}
		}
		return true, nil

	} else if strings.HasPrefix(format, "table ") {

		columns := parseTableDef(format)
		outWidth := printer.Width(ctx, data, columns)

		if !force {
			termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil && 0 < termWidth && termWidth < outWidth {
				return false, nil
			}
		}

		printer.Print(ctx, data, columns)
		return true, nil

	} else {

		lines := make([]string, 0)
		outWidth := 0

		for _, v := range data {
			str := replaceSingleLineColumnData(ctx, data, v, format)
			lines = append(lines, str)
			outWidth = mathext.Max(outWidth, printer.RealStrLen(str))
		}

		if !force {
			termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil && 0 < termWidth && termWidth < outWidth {
				return false, nil
			}
		}

		for _, v := range lines {
			ctx.PrintPrimaryOutput(v)
		}
		return true, nil

	}
}
