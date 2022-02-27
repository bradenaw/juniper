package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

// Set is a tree-structured set. Sets are a collection of unique elements. Similar to Go's built-in
// map[T]struct{} but keeps elements in sorted order.
type Set[T any] struct {
	// An extra indirect here so that tree.Set behaves like a reference type like the map builtin.
	t *btree[T, struct{}]
}

// NewSet returns a Set that uses less to determine the sort order of items. If !less(a, b) &&
// !less(b, a), then a and b are considered the same item. The output of less must not change for
// any pair of items while they are in the set.
func NewSet[T any](less xsort.Less[T]) Set[T] {
	return Set[T]{
		t: newBtree[T, struct{}](less),
	}
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return s.t.size
}

// Add adds item to the set if it is not already present.
func (s Set[T]) Add(item T) {
	s.t.Put(item, struct{}{})
}

// Remove removes item from the set if it is present, and does nothing otherwise.
func (s Set[T]) Remove(item T) {
	s.t.Delete(item)
}

// Contains returns true if item is present in the set.
func (s Set[T]) Contains(item T) bool {
	return s.t.Contains(item)
}

// First returns the lowest item in the set according to less.
func (s Set[T]) First() T {
	item, _ := s.t.First()
	return item
}

// Last returns the highest item in the set according to less.
func (s Set[T]) Last() T {
	item, _ := s.t.Last()
	return item
}

// Iterate returns an iterator that yields the elements of the set in ascending order.
//
// The set may be safely modified during iteration and the iterator will continue from the
// next-lowest item. Thus the iterator will see new items that are after the current position
// of the iterator according to less, but will not necessarily see a consistent snapshot of the
// state of the set.
func (s Set[T]) Iterate() iterator.Iterator[T] {
	return s.Range(Unbounded[T](), Unbounded[T]())
}

// Range returns an iterator that yields the elements of the set between the given bounds in
// ascending order.
//
// The set may be safely modified during iteration and the iterator will continue from the
// next-lowest item. Thus the iterator will see new items that are after the current position
// of the iterator according to less, but will not necessarily see a consistent snapshot of the
// state of the set.
func (s Set[T]) Range(lower Bound[T], upper Bound[T]) iterator.Iterator[T] {
	return iterator.Map(s.t.Range(lower, upper), func(pair KVPair[T, struct{}]) T {
		return pair.Key
	})
}

// RangeReverse returns an iterator that yields the elements of the set between the given bounds in
// descending order.
//
// The set may be safely modified during iteration and the iterator will continue from the
// next-lowest item. Thus the iterator will see new items that are after the current position
// of the iterator according to less, but will not necessarily see a consistent snapshot of the
// state of the set.
func (s Set[T]) RangeReverse(lower Bound[T], upper Bound[T]) iterator.Iterator[T] {
	return iterator.Map(s.t.RangeReverse(lower, upper), func(pair KVPair[T, struct{}]) T {
		return pair.Key
	})
}
