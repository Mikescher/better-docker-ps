package printer

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	"gogs.mikescher.com/BlackForestBytes/goext/mathext"
	"gogs.mikescher.com/BlackForestBytes/goext/termext"
	"strings"
)

type ColFun = func(ctx *cli.PSContext, cont *docker.ContainerSchema) []string

func Width(ctx *cli.PSContext, data []docker.ContainerSchema, cols []ColFun) int {
	var cells = make([][]string, 0)

	if ctx.Opt.PrintHeader {
		row := make([]string, 0)

		for _, fn := range cols {
			h := fn(ctx, nil)
			row = append(row, h[0])
		}

		cells = append(cells, row)
	}

	for _, dat := range data {
		extrow := make([][]string, 0)

		maxheight := 1
		for _, fn := range cols {
			h := fn(ctx, &dat)
			extrow = append(extrow, h)
			maxheight = mathext.Max(maxheight, len(h))
		}

		for yy := 0; yy < maxheight; yy++ {
			row := make([]string, len(cols))
			for xx := 0; xx < len(cols); xx++ {
				if yy < len(extrow[xx]) {
					row[xx] = extrow[xx][yy]
				}
			}
			cells = append(cells, row)
		}

	}

	lens := make([]int, len(cells[0]))
	for _, row := range cells {
		for i, cell := range row {
			lens[i] = mathext.Max(lens[i], RealStrLen(cell))
		}
	}

	w := 0
	for _, v := range lens {
		w += v
	}

	return w + 4*(len(cols)-1)
}

func Print(ctx *cli.PSContext, data []docker.ContainerSchema, cols []ColFun) {

	var cells = make([][]string, 0)

	if ctx.Opt.PrintHeader {
		row := make([]string, 0)

		for _, fn := range cols {
			h := fn(ctx, nil)
			row = append(row, h[0])
		}

		cells = append(cells, row)
	}

	for _, dat := range data {
		extrow := make([][]string, 0)

		maxheight := 1
		for _, fn := range cols {
			h := fn(ctx, &dat)
			extrow = append(extrow, h)
			maxheight = mathext.Max(maxheight, len(h))
		}

		for yy := 0; yy < maxheight; yy++ {
			row := make([]string, len(cols))
			for xx := 0; xx < len(cols); xx++ {
				if yy < len(extrow[xx]) {
					row[xx] = extrow[xx][yy]
				}
			}
			cells = append(cells, row)
		}

	}

	lens := make([]int, len(cells[0]))
	for _, row := range cells {
		for i, cell := range row {
			lens[i] = mathext.Max(lens[i], RealStrLen(cell))
		}
	}

	for rowidx, row := range cells {

		{
			rowstr := ""
			for colidx, cell := range row {
				if colidx > 0 {
					rowstr += "    "
				}
				if colidx == len(row)-1 {
					rowstr += cell // do not pad last
				} else {
					rowstr += TermStrPadRight(cell, " ", lens[colidx])
				}
			}
			ctx.PrintPrimaryOutput(rowstr)
		}

		if ctx.Opt.PrintHeader && ctx.Opt.PrintHeaderLines && rowidx == 0 {
			rowstr := ""
			for colidx := range row {
				if colidx > 0 {
					rowstr += "    "
				}
				rowstr += TermStrPadRight("", "-", lens[colidx])
			}
			ctx.PrintPrimaryOutput(rowstr)
		}
	}

}

func RealStrLen(cell string) int {
	return len([]rune(termext.CleanString(cell)))
}

func TermStrPadRight(str string, pad string, padlen int) string {
	if pad == "" {
		pad = " "
	}

	if RealStrLen(str) >= padlen {
		return str
	}

	return str + strings.Repeat(pad, padlen-RealStrLen(str))[0:(padlen-RealStrLen(str))]
}
