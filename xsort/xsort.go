// Package xsort contains extensions to the standard library package sort.
package xsort

import (
	"cmp"
	"slices"
	"sort"

	"github.com/bradenaw/juniper/internal/heap"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xslices"
)

// Returns true if a is less than b. Must follow the same rules as sort.Interface.Less.
type Less[T any] func(a, b T) bool

// OrderedLess is an implementation of Less for cmp.Ordered types by using the < operator.
//
// Deprecated: cmp.Less is in the standard library as of Go 1.21.
func OrderedLess[T cmp.Ordered](a, b T) bool {
	return cmp.Less(a, b)
}

func LessToCompare[T any](less Less[T]) func(T, T) int {
	return func(a, b T) int {
		if less(a, b) {
			return -1
		} else if less(b, a) {
			return 1
		}
		return 0
	}
}

// Compile-time assert the types match.
var _ Less[int] = OrderedLess[int]

// Greater returns true if a > b according to less.
func Greater[T any](less Less[T], a T, b T) bool {
	return less(b, a)
}

// LessOrEqual returns true if a <= b according to less.
func LessOrEqual[T any](less Less[T], a T, b T) bool {
	// a <= b
	// !(a > b)
	// !(b < a)
	return !less(b, a)
}

// LessOrEqual returns true if a >= b according to less.
func GreaterOrEqual[T any](less Less[T], a T, b T) bool {
	// a >= b
	// !(a < b)
	return !less(a, b)
}

// Equal returns true if a == b according to less.
func Equal[T any](less Less[T], a T, b T) bool {
	return !less(a, b) && !less(b, a)
}

// Reverse returns a Less that orders elements in the opposite order of the provided less.
func Reverse[T any](less Less[T]) Less[T] {
	return func(a, b T) bool {
		return less(b, a)
	}
}

func ReverseCompare[T any](cmp func(T, T) int) func(T, T) int {
	return func(a, b T) int {
		c := cmp(a, b)
		if -c == c {
			// Because the range of int is [-2^63, 2^63-1], -math.MinInt overflows and becomes
			// math.MinInt again, also a negative number.
			return 1
		}
		return -c
	}
}

// Slice sorts x in-place using the given less function to compare items.
//
// Follows the same rules as sort.Slice.
//
// Deprecated: slices.SortFunc is in the standard library as of Go 1.21.
func Slice[T any](x []T, less Less[T]) {
	sort.Slice(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// SliceStable stably sorts x in-place using the given less function to compare items.
//
// Follows the same rules as sort.SliceStable.
//
// Deprecated: slices.SortStableFunc is in the standard library as of Go 1.21.
func SliceStable[T any](x []T, less Less[T]) {
	sort.SliceStable(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// SliceIsSorted returns true if x is in sorted order according to the given less function.
//
// Follows the same rules as sort.SliceIsSorted.
//
// Deprecated: slices.IsSortedFunc is in the standard library as of Go 1.21.
func SliceIsSorted[T any](x []T, less Less[T]) bool {
	return sort.SliceIsSorted(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// Search searches for item in x, assumed sorted according to less, and returns the index. The
// return value is the index to insert item at if it is not present (it could be len(a)).
//
// Deprecated: slices.BinarySearchFunc is in the standard library as of Go 1.21.
func Search[T any](x []T, less Less[T], item T) int {
	return sort.Search(len(x), func(i int) bool {
		return less(item, x[i]) || !less(x[i], item)
	})
}

type valueAndSource[T any] struct {
	value  T
	source int
}

type mergeIterator[T any] struct {
	in []iterator.Iterator[T]
	h  heap.Heap[valueAndSource[T]]
}

func (iter *mergeIterator[T]) Next() (T, bool) {
	if iter.h.Len() == 0 {
		var zero T
		return zero, false
	}
	item := iter.h.Pop()
	nextItem, ok := iter.in[item.source].Next()
	if ok {
		iter.h.Push(valueAndSource[T]{nextItem, item.source})
	}
	return item.value, true
}

// Merge returns an iterator that yields all items from in in sorted order.
//
// The results are undefined if the in iterators do not yield items in sorted order according to
// less.
//
// The time complexity of Next() is O(log(k)) where k is len(in).
func Merge[T any](cmp func(T, T) int, in ...iterator.Iterator[T]) iterator.Iterator[T] {
	initial := make([]valueAndSource[T], 0, len(in))
	for i := range in {
		item, ok := in[i].Next()
		if !ok {
			continue
		}
		initial = append(initial, valueAndSource[T]{item, i})
	}
	h := heap.New(
		func(a, b valueAndSource[T]) int {
			return cmp(a.value, b.value)
		},
		func(a valueAndSource[T], i int) {},
		initial,
	)
	return &mergeIterator[T]{
		in: in,
		h:  h,
	}
}

// MergeSlices merges the already-sorted slices of in. Optionally, a pre-allocated out slice can be
// provided to store the result into.
//
// The results are undefined if the in slices are not already sorted.
//
// The time complexity is O(n * log(k)) where n is the total number of items and k is len(in).
func MergeSlices[T any](cmp func(T, T) int, out []T, in ...[]T) []T {
	n := 0
	for i := range in {
		n += len(in[i])
	}
	out = slices.Grow(out[:0], n)
	iter := Merge(cmp, xslices.Map(in, iterator.Slice[T])...)
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		out = append(out, item)
	}
	return out
}

// MinK returns the k minimum items according to cmp from iter in sorted order. If iter yields
// fewer than k items, MinK returns all of them.
func MinK[T any](cmp func(T, T) int, iter iterator.Iterator[T], k int) []T {
	h := heap.New[T](ReverseCompare(cmp), func(a T, i int) {}, nil)
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		h.Push(item)
		if h.Len() > k {
			h.Pop()
		}
	}
	out := make([]T, h.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = h.Pop()
	}
	return out
}
