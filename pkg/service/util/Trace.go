package util

import (
	"fmt"
	"runtime"
)

// Trace - debug method, will return a filename, line and caller function name
func Trace(err error) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return fmt.Sprintf(
		"=======TRACE=======\nLine: %d\nFile: %s,\nMethod: %s\nError: %s\n",
		frame.Line,
		frame.File,
		frame.Function,
		err.Error(),
	)
}
