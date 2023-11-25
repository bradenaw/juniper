//go:build !go1.21

package xmath

import "golang.org/x/exp/constraints"

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

// Clamp clamps the value of x to within min and max.
func Clamp[T constraints.Ordered](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
