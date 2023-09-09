// Package xmath contains extensions to the standard library package math.
package xmath

import "cmp"

// Min returns the minimum of a and b based on the < operator.
//
// Deprecated: min is a builtin as of Go 1.21.
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of a and b based on the > operator.
//
// Deprecated: max is a builtin as of Go 1.21.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of x. It panics if this value is not representable, for example
// because -math.MinInt32 requires more than 32 bits to represent and so does not fit in an int32.
func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64](x T) T {
	if x < 0 {
		if -x == x {
			panic("can't xmath.Abs minimum value: positive equivalent not representable")
		}
		return -x
	}
	return x
}

// Clamp clamps the value of x to within min and max.
func Clamp[T cmp.Ordered](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
