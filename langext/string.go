package langext

import (
	"fmt"
	"strings"
)

func StrPadLeft(str string, pad string, padlen int) string {
	if pad == "" {
		pad = " "
	}

	if len(str) >= padlen {
		return str
	}

	return strings.Repeat(pad, padlen-len(str))[0:(padlen-len(str))] + str
}

func StrPadRight(str string, pad string, padlen int) string {
	if pad == "" {
		pad = " "
	}

	if len(str) >= padlen {
		return str
	}

	return str + strings.Repeat(pad, padlen-len(str))[0:(padlen-len(str))]
}

func StrRunePadLeft(str string, pad string, padlen int) string {
	if pad == "" {
		pad = " "
	}

	if len([]rune(str)) >= padlen {
		return str
	}

	return strings.Repeat(pad, padlen-len([]rune(str)))[0:(padlen-len([]rune(str)))] + str
}

func StrRunePadRight(str string, pad string, padlen int) string {
	if pad == "" {
		pad = " "
	}

	if len([]rune(str)) >= padlen {
		return str
	}

	return str + strings.Repeat(pad, padlen-len([]rune(str)))[0:(padlen-len([]rune(str)))]
}

func Indent(str string, pad string) string {
	eonl := strings.HasSuffix(str, "\n")
	r := ""
	for _, v := range strings.Split(str, "\n") {
		r += pad + v + "\n"
	}

	if eonl {
		r = r[0 : len(r)-1]
	}

	return r
}

func NumToStringOpt[V IntConstraint](v *V, fallback string) string {
	if v == nil {
		return fallback
	} else {
		return fmt.Sprintf("%d", v)
	}
}
