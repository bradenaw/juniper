//go:build !go1.21

package xsort

import (
	"golang.org/x/exp/constraints"
)

// OrderedLess is an implementation of Less for cmp.Ordered types by using the < operator.
func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}

func Compare[T constraints.Ordered](a, b T) int {
	if a < b {
		return -1
	} else if b > a {
		return 1
	} else {
		return 0
	}
}
