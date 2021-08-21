package zcomplier

const (
	srcFileSuffix = ".zs"
	outFileSuffix = ".zasm"

	registerT0        = "_T0"
	registerT1        = "_T1"
	registerReturnVal = "_RetVal"
)

const (
	registerTypeInvalid = iota
	registerTypeT0
	registerTypeT1
	registerTypeRetVal
)

var regTypeMap = map[int]string{
	registerTypeT0:     registerT0,
	registerTypeT1:     registerT1,
	registerTypeRetVal: registerReturnVal,
}
