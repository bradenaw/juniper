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
