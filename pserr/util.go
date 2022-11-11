package pserr

import (
	"fmt"
	"github.com/joomcode/errorx"
)

func GetDirectOutput(err error) *errorx.Error {

	sub := err
	for sub != nil {

		errx := errorx.Cast(sub)
		if errx == nil {
			break
		}

		if uw := errx.Unwrap(); uw != nil {
			sub = uw
			continue
		}

		if errx.Type() == DirectOutput {
			return errx
		}

		sub = errx.Cause()
	}

	return nil
}

func FormatError(err error, verbose bool) string {
	if errx := GetDirectOutput(err); errx != nil {
		if verbose {
			return fmt.Sprintf("%s\n\n%+v", errx.Message(), err)
		} else {
			return errx.Message()
		}
	}

	if verbose {
		return fmt.Sprintf("%+v", err)
	} else {
		return err.Error()
	}
}

func GetExitCode(err error, fallback int) int {
	if errx := GetDirectOutput(err); errx != nil {
		if ec, ok := errx.Property(Exitcode); ok {
			if eci, ok := ec.(int); ok {
				return eci
			}
		}
	}

	return fallback
}
