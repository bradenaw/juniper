//go:build go1.21

package xsort

import (
	"cmp"
)

// OrderedLess is an implementation of Less for cmp.Ordered types by using the < operator.
//
// Deprecated: cmp.Less is in the standard library as of Go 1.21.
func OrderedLess[T cmp.Ordered](a, b T) bool {
	return cmp.Less(a, b)
}
