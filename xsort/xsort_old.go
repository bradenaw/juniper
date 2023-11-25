//go:build !go1.21

package xsort

import (
	"golang.org/x/exp/constraints"
)

// OrderedLess is an implementation of Less for cmp.Ordered types by using the < operator.
func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}
