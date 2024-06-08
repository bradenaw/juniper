//go:build !go1.21

package xsort

import (
	"golang.org/x/exp/constraints"
)

// OrderedLess is an implementation of Less for cmp.Ordered types by using the < operator.
func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}

func Compare[T constraints.Ordered](x, y T) int {
	// Copied from the standard library, here for versions older than 1.21 when it was added.
	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/cmp/cmp.go;l=40

	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}

// Copied from the standard library, here for versions older than 1.21 when it was added.
// https://cs.opensource.google/go/go/+/refs/tags/go1.22.4:src/cmp/cmp.go;l=40
//
// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}
