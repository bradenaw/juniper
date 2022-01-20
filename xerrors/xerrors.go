package xerrors

import (
	"runtime"
	"strconv"
	"strings"
)

type withStack struct {
	inner error
	pc    []uintptr
}

func (err withStack) Error() string {
	var sb strings.Builder
	frames := runtime.CallersFrames(err.pc)
	_, _ = sb.WriteString(err.inner.Error())
	_, _ = sb.WriteString("\n\n")
	for {
		frame, more := frames.Next()

		_, _ = sb.WriteString(frame.Function)
		_, _ = sb.WriteString("(...)\n        ")
		_, _ = sb.WriteString(frame.File)
		_, _ = sb.WriteString(":")
		_, _ = sb.WriteString(strconv.Itoa(frame.Line))
		_, _ = sb.WriteString("\n")

		if !more {
			break
		}
	}
	return sb.String()
}

func (err withStack) Unwrap() error {
	return err.inner
}

// WithStack returns an error that wraps err and adds the call stack of the call to WithStack to
// Error().
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	var buf [64]uintptr
	var ptrs []uintptr
	skip := 2
	for {
		n := runtime.Callers(skip, buf[:])
		ptrs = append(ptrs, buf[:n]...)
		if n < len(buf) {
			break
		}
		skip += n
	}
	return withStack{
		inner: err,
		pc:    ptrs,
	}
}
