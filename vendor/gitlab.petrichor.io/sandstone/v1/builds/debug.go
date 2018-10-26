// +build debug

package builds

import (
	"runtime"
)

const Debug bool = true

var lg logger

func SetLog(l logger) {
	lg = l
}

func DebugOut(msg string, args ...interface{}) {
	lg.Debugf(msg, args...)
}

func DebugStack(full bool) {
	b := make([]byte, 1024)
	runtime.Stack(b, full)
	lg.Debugf(string(b))
}
