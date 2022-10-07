package impl

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	pserr "better-docker-ps/fferr"
	"better-docker-ps/printer"
	"encoding/json"
)

func Execute(ctx *cli.PSContext) error {
	jsonraw, err := docker.ListContainer(ctx)
	if err != nil {
		return err
	}

	var data []docker.ContainerSchema
	err = json.Unmarshal(jsonraw, &data)
	if err != nil {
		return pserr.DirectOutput.Wrap(err, "Failed to decode Docker API response")
	}

	columns := []printer.ColFun{ //TODO make configurable
		ColContainerID,
		ColName,
		//ColFullImage,
		//ColRegistry,
		//ColImage,
		ColImageTag,
		//ColCommand,
		//ColShortCommand,
		ColCreated,
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
