// Package xerrors contains extensions to the standard library package errors.
package xerrors

import (
	"errors"
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

var noError = errors.New("no error")

// WithStack returns an error that wraps err and adds the call stack of the call to WithStack to
// Error(). If err is nil or already has a stack attached, returns err.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	// use noError in case anything along err's chain has a custom Is that calls Error() for some
	// reason.
	if errors.Is(err, withStack{inner: noError}) {
		return err
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
