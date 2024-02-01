package require2

import (
	"errors"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/exp/constraints"
)

func Equal[T comparable](t *testing.T, expected T, actual T) {
	if expected != actual {
		fatalf(t, "assertion failed: %#v == %#v", expected, actual)
	}
}

func DeepEqual[T any](t *testing.T, a T, b T) {
	if !reflect.DeepEqual(a, b) {
		fatalf(t, "assertion failed: reflect.DeepEqual(%#v, %#v)", a, b)
	}
}

func Equalf[T comparable](t *testing.T, expected T, actual T, s string, fmtArgs ...any) {
	if expected != actual {
		fatalf(t, "assertion failed: %#v == %#v\n"+s, append([]any{expected, actual}, fmtArgs...)...)
	}
}

func SlicesEqual[T comparable](t *testing.T, expected []T, actual []T) {
	n := len(expected)
	if len(actual) < n {
		n = len(actual)
	}

	for i := 0; i < n; i++ {
		if expected[i] != actual[i] {
			fatalf(t, "differ at index %d: %#v != %#v", i, expected[i], actual[i])
		}
	}

	if len(expected) != len(actual) {
		fatalf(t, "lengths differ: %d != %d", len(expected), len(actual))
	}
}

func Nil[T any](t *testing.T, a *T) {
	if a != nil {
		t.Fatal("expected nil")
	}
}

func NotNil[T any](t *testing.T, a *T) {
	if a == nil {
		t.Fatal("expected not nil")
	}
}

func NotNilf[T any](t *testing.T, a *T, f string, fmtArgs ...any) {
	if a == nil {
		fatalf(t, "expected not nil\n"+f, fmtArgs...)
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		fatalf(t, "expected no error, got %#v", err)
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		fatalf(t, "expected %T (%s), got no error", err, err)
	}
}

func ErrorIs(t *testing.T, err error, match error) {
	if !errors.Is(err, match) {
		fatalf(t, "expected %T (%s), got %#v", match, match, err)
	}
}

func Greater[T constraints.Ordered](t *testing.T, a T, b T) {
	if !(a > b) {
		fatalf(t, "assertion failed: %#v > %#v", a, b)
	}
}

func GreaterOrEqual[T constraints.Ordered](t *testing.T, a T, b T) {
	if !(a >= b) {
		fatalf(t, "assertion failed: %#v >= %#v", a, b)
	}
}

func Less[T constraints.Ordered](t *testing.T, a T, b T) {
	if !(a < b) {
		fatalf(t, "assertion failed: %#v < %#v", a, b)
	}
}

func LessOrEqual[T constraints.Ordered](t *testing.T, a T, b T) {
	if !(a <= b) {
		fatalf(t, "assertion failed: %#v <= %#v", a, b)
	}
}

func InDelta[T ~float32 | ~float64](t *testing.T, actual T, expected T, delta T) {
	diff := actual - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > delta {
		fatalf(t, "expected %#v to be within %#v of %#v, actually %#v", actual, delta, expected, diff)
	}
}

func True(t *testing.T, b bool) {
	if !b {
		fatalf(t, "expected true")
	}
}

func Truef(t *testing.T, b bool, s string, fmtArgs ...any) {
	if !b {
		fatalf(t, "expected true\n"+s, fmtArgs...)
	}
}

func ElementsMatch[T comparable](t *testing.T, a []T, b []T) {
	aSet := make(map[T]struct{}, len(a))
	for _, ai := range a {
		aSet[ai] = struct{}{}
	}
	bSet := make(map[T]struct{}, len(b))
	for _, bi := range b {
		bSet[bi] = struct{}{}
	}

	for ai := range aSet {
		_, ok := bSet[ai]
		if !ok {
			fatalf(t, "%#v appears in a but not in b", ai)
		}
	}
	for bi := range bSet {
		_, ok := aSet[bi]
		if !ok {
			fatalf(t, "%#v appears in b but not in a", bi)
		}
	}
}

func ElementsEqual[T comparable](t *testing.T, a []T, b []T) {
	aSet := make(map[T]int, len(a))
	for _, ai := range a {
		aSet[ai] += 1
	}
	bSet := make(map[T]int, len(b))
	for _, bi := range b {
		bSet[bi] += 1
	}

	for elem, aCount := range aSet {
		bCount, _ := bSet[elem]
		if aCount != bCount {
			fatalf(t, "%#v appears %d times in a and %d times in b", elem, aCount, bCount)
		}
	}
	for elem, bCount := range bSet {
		aCount, _ := aSet[elem]
		if bCount != aCount {
			fatalf(t, "%#v appears %d times in b and %d times in a", elem, bCount, aCount)
		}
	}
}

func fatalf(t *testing.T, s string, fmtArgs ...any) {
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
	var sb strings.Builder
	frames := runtime.CallersFrames(ptrs)
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

	t.Fatalf(s+"\n\n%s", append(fmtArgs, sb.String())...)
}
