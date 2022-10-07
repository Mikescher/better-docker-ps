package pserr

import "github.com/joomcode/errorx"

var (
	DopsErrors = errorx.NewNamespace("dops")
)

var (
	DirectOutput = DopsErrors.NewType("direct_out")
)

var (
	Exitcode = errorx.RegisterProperty("ffsync.exitcode")
)
