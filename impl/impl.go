package impl

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	pserr "better-docker-ps/fferr"
	"better-docker-ps/printer"
	"encoding/json"
	"strings"
)

func Execute(ctx *cli.PSContext) error {
	if strings.Contains(ctx.Opt.Format, "{{.Size}}") {
		ctx.Opt.WithSize = true
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

	if ctx.Opt.OnlyIDs {
		for _, v := range data {
			if ctx.Opt.Truncate {
				ctx.PrintPrimaryOutput(v.ID[0:12])
			} else {
				ctx.PrintPrimaryOutput(v.ID)
			}
		}
		return nil
	} else if strings.HasPrefix(ctx.Opt.Format, "table ") {

		//TODO

	} else {

		//TODO

	}

	//TODO make configurable (--format?)
	//TODO default == auto (columns have priority and get removed based on term width ??)
	columns := []printer.ColFun{
		ColContainerID,
		ColName,
		//ColFullImage,
		//ColRegistry,
		//ColImage,
		ColImageTag,
		//ColCommand,
		//ColShortCommand,
		ColRunningFor,
		ColState,
		ColStatus,
		ColPorts,
		//ColSize,
		//ColMounts,
		ColIP,
	}

	printer.Print(ctx, data, columns)

	return nil
}
