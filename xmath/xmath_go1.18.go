//go:build go1.18

// Package xmath contains extensions to the standard library package math.
package xmath

import "constraints"

// Min returns the minimum of a and b based on the < operator.
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of a and b based on the > operator.
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of x. It panics if this value is not representable, for example
// because -math.MinInt32 requires more than 32 bits to represent and so does not fit in an int32.
func Abs[T constraints.Signed](x T) T {
	if x < 0 {
		if -x == x {
			panic("can't xmath.Abs minimum value: positive equivalent not representable")
		}
		return -x
	}
	return x
}
