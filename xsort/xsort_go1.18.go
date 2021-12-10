//go:build go1.18

// Package xsort contains extensions to the standard library package sort.
package xsort

import (
	"constraints"
	"sort"

	"github.com/bradenaw/juniper/internal/heap"
	"github.com/bradenaw/juniper/slices"
)

// Returns true if a is less than b. Must follow the same rules as sort.Interface.Less.
type Less[T any] func(a, b T) bool

// OrderedLess is an implementation of Less for constraints.Ordered types by using the < operator.
func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}

// Reverse returns a Less that orders elements in the opposite order of the provided less.
func Reverse[T any](less Less[T]) Less[T] {
	return func(a, b T) bool {
		return !less(a, b)
	}
}

// Compile-time assert the types match.
var _ Less[int] = OrderedLess[int]

// Slice sorts x in-place using the given less function to compare items.
//
// Follows the same rules as sort.Slice.
func Slice[T any](x []T, less Less[T]) {
	sort.Slice(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// SliceStable stably sorts x in-place using the given less function to compare items.
//
// Follows the same rules as sort.SliceStable.
func SliceStable[T any](x []T, less Less[T]) {
	sort.SliceStable(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// SliceIsSorted returns true if x is in sorted order according to the given less function.
//
// Follows the same rules as sort.SliceIsSorted.
func SliceIsSorted[T any](x []T, less Less[T]) bool {
	return sort.SliceIsSorted(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

type valueAndSource[T any] struct {
	value  T
	source int
}

// Merge merges the already-sorted slices of in. Optionally, a pre-allocated out slice can be
// provided to store the result into.
//
// The results are undefined if the in slices are not already sorted.
//
// The complexity is O(n * log(k)) where n is the total number of items and k is len(in).
func Merge[T any](less Less[T], out []T, in ...[]T) []T {
	initial := make([]valueAndSource[T], 0, len(in))
	n := 0
	for j := range in {
		n += len(in[j])
		if len(in[j]) > 0 {
			initial = append(initial, valueAndSource[T]{in[j][0], j})
			in[j] = in[j][1:]
		}
	}
	h := heap.New(
		func(a, b valueAndSource[T]) bool {
			return less(a.value, b.value)
		},
		func(a valueAndSource[T], i int) {},
		initial,
	)
	out = slices.Grow(out[:0], n)
	for h.Len() > 0 {
		item := h.Pop()
		out = append(out, item.value)
		if len(in[item.source]) > 0 {
			h.Push(valueAndSource[T]{in[item.source][0], item.source})
			in[item.source] = in[item.source][1:]
		}
	}
	return out
}
