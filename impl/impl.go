package impl

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	"better-docker-ps/printer"
	"better-docker-ps/pserr"
	"encoding/json"
	"os"
	"strings"

	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"golang.org/x/term"
)

func Execute(ctx *cli.PSContext) error {
	for _, fmt := range ctx.Opt.Format {
		if strings.Contains(fmt, "{{.Size}}") {
			ctx.Opt.WithSize = true
		}
	}

	jsonraw, err := docker.ListContainer(ctx)
	if err != nil {
		return err
	}

	var data []docker.ContainerSchema
	err = json.Unmarshal(jsonraw, &data)
	if err != nil {
		return pserr.DirectOutput.Wrap(err, "Failed to decode Docker API response")
	}

	for i, v := range ctx.Opt.Format {

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
			str := replaceSingleLineColumnData(ctx, v, format)
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
