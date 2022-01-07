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

// Ordering defines a total ordering of T.
//
// Only the zero-value of Orderings are used. Thus Orderings usually have an underlying type of
// struct{}.
type Ordering[T any] interface {
	// Less returns true if a is less than b. Must follow the same rules as sort.Interface.Less.
	Less(a, b T) bool
}

// NaturalOrder is the Ordering of the < operator.
type NaturalOrder[T constraints.Ordered] struct{}

func (NaturalOrder[T]) Less(a, b T) bool {
	return a < b
}

// Compile-time assert the types match.
var _ Ordering[int] = NaturalOrder[int]{}

// Greater returns true if a > b according to O.
func Greater[O Ordering[T], T any](a T, b T) bool {
	var ordering O
	return ordering.Less(b, a)
}

// LessOrEqual returns true if a <= b according to O.
func LessOrEqual[O Ordering[T], T any](a T, b T) bool {
	var ordering O
	// a <= b
	// !(a > b)
	// !(b < a)
	return !ordering.Less(b, a)
}

// GreaterOrEqual returns true if a >= b according to O.
func GreaterOrEqual[O Ordering[T], T any](a T, b T) bool {
	var ordering O
	// a >= b
	// !(a < b)
	return !ordering.Less(a, b)
}

// Equal returns true if a == b according to O.
func Equal[O Ordering[T], T any](a T, b T) bool {
	var ordering O
	return !ordering.Less(a, b) && !ordering.Less(b, a)
}

// Reverse returns an Ordering that orders elements in the opposite order of the provided less.
type Reverse[O Ordering[T], T any] struct{}

func (Reverse[O, T]) Less(a, b T) bool {
	var ordering O
	return ordering.Less(b, a)
}

// Less returns true if a is less than b. Must follow the same rules as sort.Interface.Less.
type Less[T any] func(a, b T) bool

// ReverseFunc returns a Less that orders elements in the opposite order as the provided less.
func ReverseFunc[T any](less Less[T]) Less[T] {
	return func(a, b T) bool {
		return less(b, a)
	}
}

// GreaterFunc returns true if a > b according to less.
func GreaterFunc[T any](less Less[T], a T, b T) bool {
	return less(b, a)
}

// LessOrEqualFunc returns true if a <= b according to less.
func LessOrEqualFunc[T any](less Less[T], a T, b T) bool {
	// a <= b
	// !(a > b)
	// !(b < a)
	return !less(b, a)
}

// GreaterOrEqualFunc returns true if a >= b according to less.
func GreaterOrEqualFunc[T any](less Less[T], a T, b T) bool {
	// a >= b
	// !(a < b)
	return !less(a, b)
}

// EqualFunc returns true if a == b according to less.
func EqualFunc[T any](less Less[T], a T, b T) bool {
	return !less(a, b) && !less(b, a)
}

// Slice sorts x in-place using the ordering O.
//
// Follows the same rules as sort.Slice.
func Slice[O Ordering[T], T any](x []T) {
	var ordering O
	sort.Slice(x, func(i, j int) bool {
		return ordering.Less(x[i], x[j])
	})
}

// SliceFunc sorts x in-place using the given less function.
//
// Follows the same rules as sort.Slice.
func SliceFunc[T any](x []T, less Less[T]) {
	sort.Slice(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// SliceStable stably sorts x in-place using O to compare items.
//
// Follows the same rules as sort.SliceStable.
func SliceStable[O Ordering[T], T any](x []T) {
	var ordering O
	sort.SliceStable(x, func(i, j int) bool {
		return ordering.Less(x[i], x[j])
	})
}

// SliceStableFunc stably sorts x in-place using the given less function.
//
// Follows the same rules as sort.SliceStable.
func SliceStableFunc[T any](x []T, less Less[T]) {
	sort.SliceStable(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// SliceIsSorted returns true if x is in sorted order according to O.
//
// Follows the same rules as sort.SliceIsSorted.
func SliceIsSorted[O Ordering[T], T any](x []T) bool {
	var ordering O
	return sort.SliceIsSorted(x, func(i, j int) bool {
		return ordering.Less(x[i], x[j])
	})
}

// SliceIsSortedFunc returns true if x is in sorted order according to less.
//
// Follows the same rules as sort.SliceIsSorted.
func SliceIsSortedFunc[T any](x []T, less Less[T]) {
	sort.SliceStable(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

// Search searches for item in x, assumed sorted according to O, and returns the index. The return
// value is the index to insert item at if it is not present (it could be len(x)).
func Search[O Ordering[T], T any](x []T, item T) int {
	return sort.Search(len(x), func(i int) bool {
		return LessOrEqual[O](item, x[i])
	})
}

// SearchFunc searches for item in x, assumed sorted according to less, and returns the index. The
// return value is the index to insert item at if it is not present (it could be len(x)).
func SearchFunc[T any](x []T, item T, less Less[T]) int {
	return sort.Search(len(x), func(i int) bool {
		return less(item, x[i]) || !less(x[i], item)
	})
}

type valueAndSource[T any] struct {
	value  T
	source int
}

type valueAndSourceOrdering[O Ordering[T], T any] struct{}

func (valueAndSourceOrdering[O, T]) Less(a, b valueAndSource[T]) bool {
	var valueOrdering O
	return valueOrdering.Less(a.value, b.value)
}

type mergeIterator[O Ordering[T], T any] struct {
	in []iterator.Iterator[T]
	h  heap.Heap[valueAndSourceOrdering[O, T], valueAndSource[T]]
}

func (iter *mergeIterator[O, T]) Next() (T, bool) {
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
// O.
//
// The time complexity of Next() is O(log(k)) where k is len(in).
func Merge[O Ordering[T], T any](in ...iterator.Iterator[T]) iterator.Iterator[T] {
	initial := make([]valueAndSource[T], 0, len(in))
	for i := range in {
		item, ok := in[i].Next()
		if !ok {
			continue
		}
		initial = append(initial, valueAndSource[T]{item, i})
	}
	h := heap.New[valueAndSourceOrdering[O, T]](
		func(a valueAndSource[T], i int) {},
		initial,
	)
	return &mergeIterator[O, T]{
		in: in,
		h:  h,
	}
}

// Merge merges the already-sorted slices of in. Optionally, a pre-allocated out slice can be
// provided to store the result into.
//
// The results are undefined if the in slices are not already sorted.
//
// The time complexity is O(n * log(k)) where n is the total number of items and k is len(in).
func MergeSlices[O Ordering[T], T any](out []T, in ...[]T) []T {
	n := 0
	for i := range in {
		n += len(in[i])
	}
	out = slices.Grow(out[:0], n)
	inIters := make([]iterator.Iterator[T], len(in))
	for i := range in {
		inIters[i] = iterator.Slice(in[i])
	}
	iter := Merge[O](inIters...)
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		out = append(out, item)
	}
	return out
}
