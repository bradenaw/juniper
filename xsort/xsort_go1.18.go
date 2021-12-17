//go:build go1.18

// Package xsort contains extensions to the standard library package sort.
package xsort

import (
	"constraints"
	"sort"

	"github.com/bradenaw/juniper/internal/heap"
	"github.com/bradenaw/juniper/iterator"
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
		return less(b, a)
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

// Search searches for item in x, assumed sorted according to less, and returns the index. The
// return value is the index to insert item at if it is not present (it could be len(a)).
func Search[T any](x []T, less Less[T], item T) int {
	return sort.Search(len(x), func(i int) bool {
		return less(item, x[i]) || !less(x[i], item)
	})
}

type valueAndSource[T any] struct {
	value  T
	source int
}

// Merge returns an iterator that yields all items from in in sorted order.
//
// The results are undefined if the in iterators do not yield items in sorted order according to
// less.
//
// The time complexity of Next() is O(log(k)) where k is len(in).
func Merge[T any](less Less[T], in ...iterator.Iterator[T]) iterator.Iterator[T] {
	initial := make([]valueAndSource[T], 0, len(in))
	for i := range in {
		if in[i].Next() {
			initial = append(initial, valueAndSource[T]{in[i].Item(), i})
		}
	}
	h := heap.New(
		func(a, b valueAndSource[T]) bool {
			return less(a.value, b.value)
		},
		func(a valueAndSource[T], i int) {},
		initial,
	)
	return iterator.New(func() (T, bool) {
		if h.Len() == 0 {
			var zero T
			return zero, false
		}
		item := h.Pop()
		if in[item.source].Next() {
			h.Push(valueAndSource[T]{in[item.source].Item(), item.source})
		}
		return item.value, true
	})
}

// Merge merges the already-sorted slices of in. Optionally, a pre-allocated out slice can be
// provided to store the result into.
//
// The results are undefined if the in slices are not already sorted.
//
// The time complexity is O(n * log(k)) where n is the total number of items and k is len(in).
func MergeSlices[T any](less Less[T], out []T, in ...[]T) []T {
	n := 0
	for i := range in {
		n += len(in[i])
	}
	out = slices.Grow(out[:0], n)
	inIters := make([]iterator.Iterator[T], len(in))
	for i := range in {
		inIters[i] = iterator.Slice(in[i])
	}
	iter := Merge(less, inIters...)
	for iter.Next() {
		out = append(out, iter.Item())
	}
	return out
}
