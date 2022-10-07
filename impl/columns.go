package impl

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	"better-docker-ps/langext"
	"better-docker-ps/langext/term"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ColContainerID(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"CONTAINER ID"}
	}

	return []string{cont.ID[0:12]}
}

func ColFullImage(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"IMAGE"}
	}

	return []string{cont.Image}
}

func ColRegistry(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"REGISTRY"}
	}

	v, _, _ := docker.SplitDockerImage(cont.Image)

	return []string{v}
}

func ColImage(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"IMAGE"}
	}

	_, v, _ := docker.SplitDockerImage(cont.Image)

	return []string{v}
}

func ColImageTag(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"TAG"}
	}

	_, _, v := docker.SplitDockerImage(cont.Image)

	return []string{v}
}

func ColCommand(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"COMMAND"}
	}

	return []string{cont.Command}
}

func ColShortCommand(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"COMMAND"}
	}

	spl := strings.Split(cont.Command, " ")
	if len(spl) == 0 {
		return []string{""}
	} else {
		return []string{spl[0]}
	}

}

func ColCreated(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"CREATED"}
	}

	ts := time.Unix(cont.Created, 0)
	diff := time.Now().Sub(ts)

	return []string{langext.FormatNaturalDurationEnglish(diff)}
}

func ColState(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"STATE"}
	}

	strstate := "[" + strings.ToUpper(string(cont.State)) + "]"

	if !ctx.Opt.OutputColor {
		return []string{strstate}
	}

	return []string{stateColor(cont.State, strstate)}
}

func ColStatus(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"STATUS"}
	}

	if !ctx.Opt.OutputColor {
		return []string{cont.Status}
	}

	return []string{stateColor(cont.State, cont.Status)}
}

func ColPorts(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"PORTS"}
	}

	r := make(map[string]bool)
	for _, port := range cont.Ports {
		p1 := langext.StrPadLeft(strconv.Itoa(port.PublicPort), " ", 5)
		p2 := langext.StrPadLeft(strconv.Itoa(port.PrivatePort), " ", 5)

		if port.PublicPort == 0 {
			r[fmt.Sprintf("         %s / %s", p2, port.Type)] = true
		} else {
			r[fmt.Sprintf("%s -> %s / %s", p1, p2, port.Type)] = true
		}
	}

	return langext.MapKeyArr(r)
}

func ColName(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"NAME"}
	}

	r := make([]string, 0, len(cont.Names))
	for _, n := range cont.Names {
		if len(n) > 0 && n[0] == '/' {
			n = n[1:]
		}
		r = append(r, n)
	}

	return r
}

func ColSize(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"SIZE"}
	}

	if cont.SizeRw == 0 && cont.SizeRootFs == 0 {
		return []string{}
	}

	return []string{fmt.Sprintf("%v (virt %v)", langext.StrPadRight(langext.FormatBytes(cont.SizeRw), " ", 11), langext.FormatBytes(cont.SizeRootFs))}
}

func ColMounts(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"MOUNTS"}
	}

	r := make([]string, 0, len(cont.Mounts))
	for _, mnt := range cont.Mounts {
		val := fmt.Sprintf("%s -> %s", mnt.Source, mnt.Destination)
		if !mnt.RW {
			val += " [ro]"
		}
		r = append(r, val)
	}

	return r
}

func ColIP(ctx *cli.PSContext, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"IP"}
	}

	r := make([]string, 0, len(cont.NetworkSettings.Networks))
	for _, nw := range cont.NetworkSettings.Networks {
		if nw.IPAddress != "" {
			r = append(r, nw.IPAddress)
		}
	}

	return r
}

func stateColor(state docker.ContainerState, value string) string {
	switch state {
	case docker.StateCreated:
		return term.Yellow(value)
	case docker.StateRunning:
		return term.Green(value)
	case docker.StateRestarting:
		return term.Yellow(value)
	case docker.StateExited:
		return term.Red(value)
	case docker.StatePaused:
		return term.Yellow(value)
	case docker.StateDead:
		return term.Red(value)
	}
	return value
}
