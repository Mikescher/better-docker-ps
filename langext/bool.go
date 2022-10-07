package langext

func FormatBool(v bool, strTrue string, strFalse string) string {
	if v {
		return strTrue
	} else {
		return strFalse
	}
}
