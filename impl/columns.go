package impl

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	"better-docker-ps/printer"
	"bytes"
	"encoding/json"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/rext"
	"gogs.mikescher.com/BlackForestBytes/goext/termext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var rexIP = rext.W(regexp.MustCompile("^(?P<b0>[0-9]{1,3})\\.(?P<b1>[0-9]{1,3})\\.(?P<b2>[0-9]{1,3})\\.(?P<b3>[0-9]{1,3})$"))

type ColSortFun = func(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int

type ColumnDef struct {
	Reader printer.ColFun
	Sorter ColSortFun
}

var ColumnMap = map[string]ColumnDef{
	"ID":                  {ColContainerID, SortContainerID},
	"Image":               {ColFullImage, SortFullImage},
	"ImageName":           {ColImage, SortImage},
	"ImageTag":            {ColImageTag, SortImageTag},
	"Registry":            {ColRegistry, SortRegistry},
	"ImageRegistry":       {ColRegistry, SortRegistry},
	"Tag":                 {ColImageTag, SortImageTag},
	"Command":             {ColCommand, SortCommand},
	"ShortCommand":        {ColShortCommand, SortShortCommand},
	"CreatedAt":           {ColCreatedAt, SortCreatedAt},
	"RunningFor":          {ColRunningFor, SortRunningFor},
	"Ports":               {ColPortsPublished, SortPortsPublished},
	"PublishedPorts":      {ColPortsPublished, SortPortsPublished},
	"ShortPublishedPorts": {ColPortsPublishedShort, SortPortsPublishedShort},
	"LongPublishedPorts":  {ColPortsPublishedLong, SortPortsPublishedLong},
	"ExposedPorts":        {ColPortsExposed, SortPortsExposed},
	"NotPublishedPorts":   {ColPortsNotPublished, SortPortsNotPublished},
	"PublicPorts":         {ColPortsPublicPart, SortPortsPublicPart},
	"State":               {ColState, SortState},
	"Status":              {ColStatus, SortStatus},
	"Size":                {ColSize, SortSize},
	"Names":               {ColName, SortName},
	"Labels":              {ColLabels, SortLabels},
	"LabelKeys":           {ColLabelKeys, SortLabelKeys},
	"Mounts":              {ColMounts, SortMounts},
	"Networks":            {ColNetworks, SortNetworks},
	"IP":                  {ColIP, SortIP},
}

func ColContainerID(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"CONTAINER ID"}
	}

	if ctx.Opt.Truncate {
		return []string{cont.ID[0:12]}
	} else {
		return []string{cont.ID}
	}
}

func ColFullImage(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"IMAGE"}
	}

	return []string{cont.Image}
}

func ColRegistry(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"REGISTRY"}
	}

	v, _, _ := docker.SplitDockerImage(ctx, cont.Image)

	return []string{v}
}

func ColImage(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"IMAGE"}
	}

	_, v, _ := docker.SplitDockerImage(ctx, cont.Image)

	return []string{v}
}

func ColImageTag(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"TAG"}
	}

	_, _, v := docker.SplitDockerImage(ctx, cont.Image)

	return []string{v}
}

func ColCommand(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"COMMAND"}
	}

	cmd := cont.Command
	if ctx.Opt.Truncate && len(cmd) > 20 {
		cmd = cmd[:19] + "…"
	}

	cmd = "\"" + cmd + "\""

	return []string{cmd}
}

func ColShortCommand(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
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

func ColRunningFor(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"CREATED"}
	}

	ts := time.Unix(cont.Created, 0)
	diff := time.Now().Sub(ts)

	return []string{timeext.FormatNaturalDurationEnglish(diff)}
}

func ColCreatedAt(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		if ctx.Opt.TimeFormatHeader != "" {
			hdr := time.Now().In(ctx.Opt.TimeZone).Format(ctx.Opt.TimeFormatHeader)
			if hdr == "Z UTC" {
				hdr = "UTC"
			}
			return []string{"CREATED AT (" + hdr + ")"}
		} else {
			return []string{"CREATED AT"}
		}
	}

	ts := time.Unix(cont.Created, 0)

	return []string{ts.In(ctx.Opt.TimeZone).Format(ctx.Opt.TimeFormat)}
}

func ColState(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"STATE"}
	}

	strstate := "[" + strings.ToUpper(string(cont.State)) + "]"

	if !ctx.Opt.OutputColor {
		return []string{strstate}
	}

	return []string{stateColor(cont.State, strstate)}
}

func ColStatus(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"STATUS"}
	}

	if !ctx.Opt.OutputColor {
		return []string{cont.Status}
	}

	return []string{statusColor(cont.Status, cont.Status)}
}

func ColPortsExposed(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"EXPOSED PORTS"}
	}

	m := make(map[string]bool)
	r := make([]string, 0)
	for _, port := range cont.PortsSorted() {
		p1 := langext.StrPadLeft(strconv.Itoa(port.PublicPort), " ", 5)
		p2 := langext.StrPadLeft(strconv.Itoa(port.PrivatePort), " ", 5)

		if port.PublicPort == 0 {
			str := fmt.Sprintf("         %s / %s", p2, port.Type)
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, str)
			}
		} else {
			str := fmt.Sprintf("%s -> %s / %s", p1, p2, port.Type)
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, str)
			}
		}
	}

	return r
}

func ColPortsPublicPart(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"EXPOSED PORTS"}
	}

	m := make(map[string]bool)
	r := make([]string, 0)
	for _, port := range cont.PortsSorted() {
		if port.PublicPort != 0 {
			str := fmt.Sprintf("%d", port.PublicPort)
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, strconv.Itoa(port.PublicPort))
			}
		}
	}

	return r
}

func ColPortsPublished(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"PUBLISHED PORTS"}
	}

	pubPortLenMax := ctx.GetIntFromCache("Printer::Ports::port_pub_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(strconv.Itoa(v2.PublicPort)))
			}
		}
		return ml
	})

	privPortLenMax := ctx.GetIntFromCache("Printer::Ports::port_pub_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(strconv.Itoa(v2.PrivatePort)))
			}
		}
		return ml
	})

	m := make(map[string]bool)
	r := make([]string, 0)
	for _, port := range cont.PortsSorted() {
		p1 := langext.StrPadLeft(strconv.Itoa(port.PublicPort), " ", pubPortLenMax)
		p2 := langext.StrPadLeft(strconv.Itoa(port.PrivatePort), " ", privPortLenMax)

		if port.PublicPort != 0 {
			str := fmt.Sprintf("%s -> %s / %s", p1, p2, port.Type)
			if port.IsLoopback() {
				str += " (loc)"
			}
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, str)
			}
		}
	}

	return r
}

func ColPortsPublishedShort(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"PUBLISHED PORTS"}
	}

	pubPortLenMax := ctx.GetIntFromCache("Printer::Ports::port_pub_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(strconv.Itoa(v2.PublicPort)))
			}
		}
		return ml
	})

	privPortLenMax := ctx.GetIntFromCache("Printer::Ports::port_pub_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(strconv.Itoa(v2.PrivatePort)))
			}
		}
		return ml
	})

	m := make(map[string]bool)
	r := make([]string, 0)
	for _, port := range cont.PortsSorted() {
		p1 := langext.StrPadLeft(strconv.Itoa(port.PublicPort), " ", pubPortLenMax)
		p2 := langext.StrPadLeft(strconv.Itoa(port.PrivatePort), " ", privPortLenMax)

		if port.PublicPort != 0 {
			str := fmt.Sprintf("%s -> %s", p1, p2)
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, str)
			}
		}
	}

	return r
}

func ColPortsPublishedLong(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"PUBLISHED PORTS"}
	}

	iplenMax := ctx.GetIntFromCache("Printer::Ports::ip_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(v2.IP))
			}
		}
		return ml
	})

	pubPortLenMax := ctx.GetIntFromCache("Printer::Ports::port_pub_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(strconv.Itoa(v2.PublicPort)))
			}
		}
		return ml
	})

	privPortLenMax := ctx.GetIntFromCache("Printer::Ports::port_pub_length", func() int {
		ml := 0
		for _, v1 := range allData {
			for _, v2 := range v1.Ports {
				ml = mathext.Max(ml, len(strconv.Itoa(v2.PrivatePort)))
			}
		}
		return ml
	})

	m := make(map[string]bool)
	r := make([]string, 0)
	for _, port := range cont.PortsSorted() {
		p0 := langext.StrPadLeft("["+port.IP+"]", " ", iplenMax+2)
		p1 := langext.StrPadLeft(strconv.Itoa(port.PublicPort), " ", pubPortLenMax)
		p2 := langext.StrPadLeft(strconv.Itoa(port.PrivatePort), " ", privPortLenMax)

		if port.PublicPort != 0 {
			str := fmt.Sprintf("%s:%s -> %s", p0, p1, p2)
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, str)
			}
		}
	}

	return r
}

func ColPortsNotPublished(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"NOT PUBLISHED PORTS"}
	}

	m := make(map[string]bool)
	r := make([]string, 0)
	for _, port := range cont.PortsSorted() {
		p2 := langext.StrPadLeft(strconv.Itoa(port.PrivatePort), " ", 5)

		if port.PublicPort == 0 {
			str := fmt.Sprintf("%s / %s", p2, port.Type)
			if _, ok := m[str]; !ok {
				m[str] = true
				r = append(r, str)
			}
		}
	}

	return r
}

func ColName(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
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

func ColSize(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"SIZE"}
	}

	if cont.SizeRw == 0 && cont.SizeRootFs == 0 {
		return []string{}
	}

	return []string{fmt.Sprintf("%v (virt %v)", langext.StrPadRight(langext.FormatBytes(cont.SizeRw), " ", 11), langext.FormatBytes(cont.SizeRootFs))}
}

func ColMounts(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
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

func ColIP(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
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

func ColLabels(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"LABELS"}
	}

	r := make([]string, 0, len(cont.Mounts))
	for k, v := range cont.Labels {
		r = append(r, fmt.Sprintf("%s := %s", k, v))
	}

	return r
}

func ColLabelKeys(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"LABELS"}
	}

	r := make([]string, 0, len(cont.Mounts))
	for k, _ := range cont.Labels {
		r = append(r, k)
	}

	return r
}

func ColNetworks(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
	if cont == nil {
		return []string{"NETWORKS"}
	}

	r := make([]string, 0, len(cont.Mounts))
	for k := range cont.NetworkSettings.Networks {
		r = append(r, k)
	}

	return r
}

func ColPlaintext(str string) printer.ColFun {
	return func(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) []string {
		return []string{str}
	}
}

// #####################################################################################################################

func SortContainerID(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	if ctx.Opt.Truncate {
		return langext.Compare(v1.ID[0:12], v2.ID[0:12])
	} else {
		return langext.Compare(v1.ID, v2.ID)
	}
}

func SortFullImage(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.Compare(v1.Image, v2.Image)
}

func SortRegistry(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	reg1, _, _ := docker.SplitDockerImage(ctx, v1.Image)
	reg2, _, _ := docker.SplitDockerImage(ctx, v2.Image)

	return langext.Compare(reg1, reg2)
}

func SortImage(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	_, img1, _ := docker.SplitDockerImage(ctx, v1.Image)
	_, img2, _ := docker.SplitDockerImage(ctx, v2.Image)

	return langext.Compare(img1, img2)
}

func SortImageTag(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	_, _, tag1 := docker.SplitDockerImage(ctx, v1.Image)
	_, _, tag2 := docker.SplitDockerImage(ctx, v2.Image)

	return langext.Compare(tag1, tag2)
}

func SortCommand(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.Compare(v1.Command, v2.Command)
}

func SortShortCommand(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	spl1 := strings.Split(v1.Command, " ")
	sc1 := ""
	if len(spl1) > 0 {
		sc1 = spl1[0]
	}

	spl2 := strings.Split(v2.Command, " ")
	sc2 := ""
	if len(spl2) > 0 {
		sc2 = spl2[0]
	}

	return langext.Compare(sc1, sc2)
}

func SortRunningFor(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.Compare(v1.Created, v2.Created) * -1 // runnign for is 'now - created', so we need to invert the sort order
}

func SortCreatedAt(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.Compare(v1.Created, v2.Created)
}

func SortState(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.Compare(v1.State.Num(), v2.State.Num())
}

func SortStatus(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.Compare(v1.Status, v2.Status)
}

func SortPortsExposed(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	parr1 := langext.ArrCopy(v1.Ports)
	parr2 := langext.ArrCopy(v2.Ports)

	pl1 := langext.ArrMap(parr1, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})
	pl2 := langext.ArrMap(parr2, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})

	langext.SortStable(pl1)
	langext.SortStable(pl2)

	pstr1 := strings.Join(pl1, "\n")
	pstr2 := strings.Join(pl2, "\n")

	return langext.Compare(pstr1, pstr2)
}

func SortPortsPublished(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	parr1 := langext.ArrCopy(v1.Ports)
	parr2 := langext.ArrCopy(v2.Ports)

	parr1 = langext.ArrFilter(parr1, func(v docker.PortSchema) bool { return v.PublicPort != 0 })
	parr2 = langext.ArrFilter(parr2, func(v docker.PortSchema) bool { return v.PublicPort != 0 })

	pl1 := langext.ArrMap(parr1, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})
	pl2 := langext.ArrMap(parr2, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})

	langext.SortStable(pl1)
	langext.SortStable(pl2)

	pstr1 := strings.Join(pl1, "\n")
	pstr2 := strings.Join(pl2, "\n")

	return langext.Compare(pstr1, pstr2)
}

func SortPortsPublishedShort(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	parr1 := langext.ArrCopy(v1.Ports)
	parr2 := langext.ArrCopy(v2.Ports)

	parr1 = langext.ArrFilter(parr1, func(v docker.PortSchema) bool { return v.PublicPort != 0 })
	parr2 = langext.ArrFilter(parr2, func(v docker.PortSchema) bool { return v.PublicPort != 0 })

	pl1 := langext.ArrMap(parr1, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})
	pl2 := langext.ArrMap(parr2, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})

	langext.SortStable(pl1)
	langext.SortStable(pl2)

	pstr1 := strings.Join(pl1, "\n")
	pstr2 := strings.Join(pl2, "\n")

	return langext.Compare(pstr1, pstr2)
}

func SortPortsPublishedLong(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	parr1 := langext.ArrCopy(v1.Ports)
	parr2 := langext.ArrCopy(v2.Ports)

	parr1 = langext.ArrFilter(parr1, func(v docker.PortSchema) bool { return v.PublicPort != 0 })
	parr2 = langext.ArrFilter(parr2, func(v docker.PortSchema) bool { return v.PublicPort != 0 })

	pl1 := langext.ArrMap(parr1, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})
	pl2 := langext.ArrMap(parr2, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})

	langext.SortStable(pl1)
	langext.SortStable(pl2)

	pstr1 := strings.Join(pl1, "\n")
	pstr2 := strings.Join(pl2, "\n")

	return langext.Compare(pstr1, pstr2)
}

func SortPortsNotPublished(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	parr1 := langext.ArrCopy(v1.Ports)
	parr2 := langext.ArrCopy(v2.Ports)

	parr1 = langext.ArrFilter(parr1, func(v docker.PortSchema) bool { return v.PublicPort == 0 })
	parr2 = langext.ArrFilter(parr2, func(v docker.PortSchema) bool { return v.PublicPort == 0 })

	pl1 := langext.ArrMap(parr1, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})
	pl2 := langext.ArrMap(parr2, func(v docker.PortSchema) string {
		return fmt.Sprintf("%s;%08d;%08d;%s", v.IP, v.PrivatePort, v.PublicPort, v.Type)
	})

	langext.SortStable(pl1)
	langext.SortStable(pl2)

	pstr1 := strings.Join(pl1, "\n")
	pstr2 := strings.Join(pl2, "\n")

	return langext.Compare(pstr1, pstr2)
}

func SortPortsPublicPart(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	parr1 := langext.ArrCopy(v1.Ports)
	parr2 := langext.ArrCopy(v2.Ports)

	parr1 = langext.ArrFilter(parr1, func(v docker.PortSchema) bool { return v.PublicPort != 0 })
	parr2 = langext.ArrFilter(parr2, func(v docker.PortSchema) bool { return v.PublicPort != 0 })

	pl1 := langext.ArrMap(parr1, func(v docker.PortSchema) string {
		return fmt.Sprintf("%08d", v.PublicPort)
	})
	pl2 := langext.ArrMap(parr2, func(v docker.PortSchema) string {
		return fmt.Sprintf("%08d", v.PublicPort)
	})

	langext.SortStable(pl1)
	langext.SortStable(pl2)

	pstr1 := strings.Join(pl1, ":")
	pstr2 := strings.Join(pl2, ":")

	return langext.Compare(pstr1, pstr2)
}

func SortName(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	names1 := langext.ArrCopy(v1.Names)
	names2 := langext.ArrCopy(v2.Names)

	langext.SortStable(names1)
	langext.SortStable(names2)

	return langext.CompareArr(names1, names2)
}

func SortSize(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	return langext.CompareArr([]int64{v1.SizeRw, v1.SizeRootFs}, []int64{v2.SizeRw, v2.SizeRootFs})
}

func SortMounts(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	mounts1 := langext.ArrMap(v1.Mounts, func(v docker.ContainerMount) string {
		return fmt.Sprintf("%s\n%s", v.Source, v.Destination)
	})
	mounts2 := langext.ArrMap(v2.Mounts, func(v docker.ContainerMount) string {
		return fmt.Sprintf("%s\n%s", v.Source, v.Destination)
	})

	langext.SortStable(mounts1)
	langext.SortStable(mounts2)

	return langext.CompareArr(mounts1, mounts2)
}

func SortIP(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	ips1 := langext.ArrMap(langext.MapToArr(v1.NetworkSettings.Networks), func(v langext.MapEntry[string, docker.ContainerSingleNetworkSettings]) string {
		return v.Value.IPAddress
	})
	ips2 := langext.ArrMap(langext.MapToArr(v2.NetworkSettings.Networks), func(v langext.MapEntry[string, docker.ContainerSingleNetworkSettings]) string {
		return v.Value.IPAddress
	})

	ips1 = langext.ArrFilter(ips1, func(v string) bool {
		return v != ""
	})
	ips2 = langext.ArrFilter(ips2, func(v string) bool {
		return v != ""
	})

	ips1 = langext.ArrMap(ips1, func(v string) string {
		return ipExpand(v)
	})
	ips2 = langext.ArrMap(ips2, func(v string) string {
		return ipExpand(v)
	})

	langext.SortStable(ips1)
	langext.SortStable(ips2)

	return langext.CompareArr(ips1, ips2)
}

func SortLabels(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	lbls1 := langext.ArrMap(langext.MapToArr(v1.Labels), func(v langext.MapEntry[string, string]) string {
		return fmt.Sprintf("%s\n%s", v.Key, v.Value)
	})
	lbls2 := langext.ArrMap(langext.MapToArr(v2.Labels), func(v langext.MapEntry[string, string]) string {
		return fmt.Sprintf("%s\t%s", v.Key, v.Value)
	})

	langext.SortStable(lbls1)
	langext.SortStable(lbls2)

	return langext.CompareArr(lbls1, lbls2)
}

func SortLabelKeys(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	lbls1 := langext.ArrMap(langext.MapToArr(v1.Labels), func(v langext.MapEntry[string, string]) string {
		return v.Key
	})
	lbls2 := langext.ArrMap(langext.MapToArr(v2.Labels), func(v langext.MapEntry[string, string]) string {
		return v.Key
	})

	langext.SortStable(lbls1)
	langext.SortStable(lbls2)

	return langext.CompareArr(lbls1, lbls2)
}

func SortNetworks(ctx *cli.PSContext, v1 *docker.ContainerSchema, v2 *docker.ContainerSchema) int {
	ntwrk1 := langext.ArrMap(langext.MapToArr(v1.NetworkSettings.Networks), func(v langext.MapEntry[string, docker.ContainerSingleNetworkSettings]) string {
		return v.Key
	})
	ntwrk2 := langext.ArrMap(langext.MapToArr(v2.NetworkSettings.Networks), func(v langext.MapEntry[string, docker.ContainerSingleNetworkSettings]) string {
		return v.Key
	})

	langext.SortStable(ntwrk1)
	langext.SortStable(ntwrk2)

	return langext.CompareArr(ntwrk1, ntwrk2)
}

// #####################################################################################################################

func getColFun(colkey string) (printer.ColFun, bool) {

	// Fast branch, simple references to columns

	for k, v := range ColumnMap {
		if "{{."+k+"}}" == colkey {
			return v.Reader, true
		}
	}

	// Slow branch, fully-featured go templates

	if strings.HasPrefix(colkey, "{{") && strings.HasSuffix(colkey, "}}") {
		return templateColFun(colkey, ""), true
	}

	if splt := strings.SplitN(colkey, ":", 2); len(splt) == 2 && strings.HasPrefix(splt[1], "{{") && strings.HasSuffix(splt[1], "}}") {
		return templateColFun(splt[1], splt[0]), true
	}

	// Fallback, nothing

	return nil, false
}

func templateColFun(fmtstr string, header string) printer.ColFun {
	return func(ctx *cli.PSContext, allData []docker.ContainerSchema, cont *docker.ContainerSchema) (res []string) {
		defer func() {
			if r := recover(); r != nil {
				ctx.PrintErrorMessage(fmt.Sprintf("Panic in template evaluation of '%s':\n%v", fmtstr, r))
				res = []string{"@ERROR"}
			}
		}()

		if cont == nil {
			return []string{header}
		}

		funcs := template.FuncMap{
			"join": strings.Join,
			"array_last": func(v any) any {
				rval := reflect.ValueOf(v)
				alen := rval.Len()
				if alen == 0 {
					return nil
				}
				return rval.Index(alen - 1).Interface()
			},
			"array_slice": func(v any, start int, end int) any {
				rval := reflect.ValueOf(v)
				alen := rval.Len()

				start = max(0, min(alen, start))
				end = max(0, min(alen, end))

				return rval.Slice(start, end).Interface()
			},
			"in_array": func(compval any, arrval any) (resp bool) {
				defer func() {
					if rec := recover(); rec != nil {
						resp = false
					}
				}()
				v := reflect.ValueOf(arrval)
				for i := 0; i < v.Len(); i++ {
					if v.Index(i).Equal(reflect.ValueOf(compval)) {
						return true
					}
				}
				return false
			},
			"json": func(obj any) string {
				v, err := json.Marshal(obj)
				if err != nil {
					panic(err)
				}
				return string(v)
			},
			"json_indent": func(obj any) string {
				v, err := json.MarshalIndent(obj, "", "  ")
				if err != nil {
					panic(err)
				}
				return string(v)
			},
			"json_pretty": func(v string) string {
				buffer := &bytes.Buffer{}
				err := json.Indent(buffer, []byte(v), "", "  ")
				if err != nil {
					return v
				} else {
					return buffer.String()
				}
			},
			"coalesce": func(val any, def any) any {
				if langext.IsNil(val) {
					return def
				} else {
					return val
				}
			},
			"to_string": func(v any) string {
				return fmt.Sprintf("%v", v)
			},
			"deref": func(vInput any) any {
				val := reflect.ValueOf(vInput)
				if val.Kind() == reflect.Ptr {
					return val.Elem().Interface()
				}
				return ""
			},
			"now": func() time.Time {
				return time.Now()
			},
			"uniqid": func() string {
				return langext.MustRawHexUUID()
			},
		}

		templ, err := template.New("col").Funcs(funcs).Parse(fmtstr)
		if err != nil {
			ctx.PrintErrorMessage(fmt.Sprintf("Error in template parsing of '%s':\n%v", fmtstr, err.Error()))
			res = []string{"@ERROR"}
		}

		bfr := &bytes.Buffer{}
		err = templ.Execute(bfr, *cont)
		if err != nil {
			ctx.PrintErrorMessage(fmt.Sprintf("Error in template evaluation of '%s':\n%v", fmtstr, err.Error()))
			res = []string{"@ERROR"}
		}

		return strings.Split(bfr.String(), "\n")
	}
}

func getSortFun(colkey string) (ColSortFun, bool) {
	if cdef, ok := ColumnMap[colkey]; ok {
		return cdef.Sorter, true
	}
	return nil, false
}

func replaceSingleLineColumnData(ctx *cli.PSContext, allData []docker.ContainerSchema, data docker.ContainerSchema, format string) string {

	r := format

	for k, v := range ColumnMap {
		r = strings.ReplaceAll(r, "{{."+k+"}}", strings.Join(v.Reader(ctx, allData, &data), " "))
	}

	return r

}

func parseTableDef(fmt string) []printer.ColFun {
	split := strings.Split(fmt[6:], "\\t")
	columns1 := make([]printer.ColFun, 0)
	for _, v := range split {
		if cf, ok := getColFun(v); ok {
			columns1 = append(columns1, cf)
		} else {
			columns1 = append(columns1, ColPlaintext(v))
		}
	}
	return columns1
}

func stateColor(state docker.ContainerState, value string) string {
	switch state {
	case docker.StateCreated:
		return termext.Yellow(value)
	case docker.StateRunning:
		return termext.Green(value)
	case docker.StateRestarting:
		return termext.Yellow(value)
	case docker.StateExited:
		return termext.Red(value)
	case docker.StatePaused:
		return termext.Yellow(value)
	case docker.StateDead:
		return termext.Red(value)
	}
	return value
}

func statusColor(status string, value string) string {
	if status == "Created" {
		return termext.Yellow(value)
	}

	if strings.HasPrefix(status, "Exited") {
		return termext.Red(value)
	}

	if strings.HasPrefix(status, "Up") {
		if strings.HasSuffix(status, "(unhealthy)") {
			return termext.Red(value)
		}
		if strings.HasSuffix(status, "(health: starting)") {
			return termext.Yellow(value)
		}

		return termext.Green(value)
	}

	return value
}

func ipExpand(ip string) string {
	if match, ok := rexIP.MatchFirst(ip); ok {
		return fmt.Sprintf("%03s.%03s.%03s.%03s",
			match.GroupByName("b0").Value(),
			match.GroupByName("b1").Value(),
			match.GroupByName("b2").Value(),
			match.GroupByName("b3").Value())
	}
	return ip
}
