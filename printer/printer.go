package printer

import (
	"better-docker-ps/cli"
	"better-docker-ps/docker"
	"better-docker-ps/langext"
	"better-docker-ps/langext/term"
	"strings"
)

type ColFun = func(ctx *cli.PSContext, cont *docker.ContainerSchema) []string

func Print(ctx *cli.PSContext, data []docker.ContainerSchema, cols []ColFun) {

	var cells = make([][]string, 0)

	{
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
			maxheight = langext.Max(maxheight, len(h))
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
			lens[i] = langext.Max(lens[i], reallen(cell))
		}
	}

	for rowidx, row := range cells {

		{
			rowstr := ""
			for colidx, cell := range row {
				if colidx > 0 {
					rowstr += "    "
				}
				rowstr += strPadRight(cell, " ", lens[colidx])
			}
			ctx.PrintPrimaryOutput(rowstr)
		}

		if rowidx == 0 {
			rowstr := ""
			for colidx := range row {
				if colidx > 0 {
					rowstr += "    "
				}
				rowstr += strPadRight("", "-", lens[colidx])
			}
			ctx.PrintPrimaryOutput(rowstr)
		}
	}

}

func reallen(cell string) int {
	return len([]rune(term.CleanString(cell)))
}

func strPadRight(str string, pad string, padlen int) string {
	if pad == "" {
		pad = " "
	}

	if reallen(str) >= padlen {
		return str
	}

	return str + strings.Repeat(pad, padlen-reallen(str))[0:(padlen-reallen(str))]
}
