package term

import "strings"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
)

func CleanString(v string) string {
	v = strings.ReplaceAll(v, colorReset, "")
	v = strings.ReplaceAll(v, colorRed, "")
	v = strings.ReplaceAll(v, colorGreen, "")
	v = strings.ReplaceAll(v, colorYellow, "")
	v = strings.ReplaceAll(v, colorBlue, "")
	v = strings.ReplaceAll(v, colorPurple, "")
	v = strings.ReplaceAll(v, colorCyan, "")
	v = strings.ReplaceAll(v, colorGray, "")
	v = strings.ReplaceAll(v, colorWhite, "")

	return v
}

func Red(v string) string {
	return colorRed + v + colorReset
}

func Green(v string) string {
	return colorGreen + v + colorReset
}

func Yellow(v string) string {
	return colorYellow + v + colorReset
}

func Blue(v string) string {
	return colorBlue + v + colorReset
}

func Purple(v string) string {
	return colorPurple + v + colorReset
}

func Cyan(v string) string {
	return colorCyan + v + colorReset
}

func Gray(v string) string {
	return colorGray + v + colorReset
}

func White(v string) string {
	return colorWhite + v + colorReset
}
