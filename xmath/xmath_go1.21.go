//go:build go1.21

package xmath

import "cmp"

// Min returns the minimum of a and b based on the < operator.
//
// Deprecated: min is a builtin as of Go 1.21.
func Min[T cmp.Ordered](a, b T) T {
	return min(a, b)
}

// Max returns the maximum of a and b based on the > operator.
//
// Deprecated: max is a builtin as of Go 1.21.
func Max[T cmp.Ordered](a, b T) T {
	return max(a, b)
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
