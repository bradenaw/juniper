//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

// Set is a tree-structured set. Sets are a collection of unique elements. Similar to Go's built-in
// map[T]struct{} but keeps elements in sorted order.
type Set[O xsort.Ordering[T], T any] struct {
	t *tree[T, struct{}, O]
}

// NewSet returns a Set that uses less to determine the sort order of items. If !less(a, b) &&
// !less(b, a), then a and b are considered the same item.
func NewSet[O xsort.Ordering[T], T any]() Set[O, T] {
	return Set[O, T]{
		t: newTree[T, struct{}, O](),
	}
}

// Len returns the number of elements in the set.
func (s Set[O, T]) Len() int {
	return s.t.size
}

// Add adds item to the set if it is not already present.
func (s Set[O, T]) Add(item T) {
	s.t.Put(item, struct{}{})
}

// Remove removes item from the set if it is present, and does nothing otherwise.
func (s Set[O, T]) Remove(item T) {
	s.t.Delete(item)
}

// Contains returns true if item is present in the set.
func (s Set[O, T]) Contains(item T) bool {
	return s.t.Contains(item)
}

// First returns the lowest item in the set according to less.
func (s Set[O, T]) First() T {
	item, _ := s.t.First()
	return item
}

// Last returns the highest item in the set according to less.
func (s Set[O, T]) Last() T {
	item, _ := s.t.Last()
	return item
}

// Iterate returns an iterator that yields the elements of the set in sorted order.
//
// The set may be safely modified during iteration and the iterator will continue from the
// next-lowest item. Thus if the set is modified, the iterator will not necessarily return all of
// the items present in the set.
func (s Set[O, T]) Iterate() iterator.Iterator[T] {
	return s.Cursor().Forward()
}

// Cursor returns a cursor into the set placed at the first item.
func (s Set[O, T]) Cursor() *SetCursor[O, T] {
	inner := s.t.Cursor()
	return &SetCursor[O, T]{
		inner: inner,
	}
}

// SetCursor is a cursor into a Set.
//
// A cursor is usable while a set is being modified. If the item the cursor is at is deleted, the
// cursor will still return the old item.
type SetCursor[O xsort.Ordering[T], T any] struct {
	inner cursor[T, struct{}, O]
}

// SeekFirst moves the cursor to the first item in the set.
//
// SeekFirst is O(log(n)).
func (c *SetCursor[O, T]) SeekFirst() { c.inner.SeekFirst() }

// SeekLast moves the cursor to the last item in the set.
//
// SeekLast is O(log(n)).
func (c *SetCursor[O, T]) SeekLast() { c.inner.SeekLast() }

// SeekLastLess moves the cursor to the item in the set just before x.
//
// SeekLastLess is O(log(n)).
func (c *SetCursor[O, T]) SeekLastLess(x T) { c.inner.SeekLastLess(x) }

// SeekLastLessOrEqual moves the cursor to the greatest item in the set that is less than or equal
// to x.
//
// SeekLastLessOrEqual is O(log(n)).
func (c *SetCursor[O, T]) SeekLastLessOrEqual(x T) { c.inner.SeekLastLessOrEqual(x) }

// SeekFirstGreaterOrEqual moves the cursor to the least item in the set that is greater than or
// equal to x.
//
// SeetFirstGreaterOrEqual is O(log(n)).
func (c *SetCursor[O, T]) SeekFirstGreaterOrEqual(x T) { c.inner.SeekFirstGreaterOrEqual(x) }

// SeekFirstGreater moves the cursor to the item in the set just after x.
//
// SeekFirstGreater is O(log(n)).
func (c *SetCursor[O, T]) SeekFirstGreater(x T) { c.inner.SeekFirstGreater(x) }

// Next moves the cursor to the next item in the set.
//
// Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
// case it's O(log(n)).
func (c *SetCursor[O, T]) Next() { c.inner.Next() }

// Prev moves the cursor to the previous item in the set.
//
// Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
// case it's O(log(n)).
func (c *SetCursor[O, T]) Prev() { c.inner.Prev() }

// Ok returns false if the cursor is not currently placed at an item, for example if Next advances
// past the last item.
func (c *SetCursor[O, T]) Ok() bool { return c.inner.Ok() }

// Item returns the item that the cursor is at. Panics if Ok is false.
func (c *SetCursor[O, T]) Item() T { return c.inner.Key() }

// Forward returns an iterator that starts from the cursor's position and yields all of the elements
// greater than or equal to the cursor in ascending order.
//
// This iterator's Next method is amoritized O(1), unless the map changes in which case the
// following Next is O(log(n)) where n is the number of elements in the map.
func (c *SetCursor[O, T]) Forward() iterator.Iterator[T] {
	return iterator.Map(c.inner.Forward(), func(kv KVPair[T, struct{}]) T {
		return kv.K
	})
}

// Backward returns an iterator that starts from the cursor's position and yields all of the
// elements less than or equal to the cursor in descending order.
//
// This iterator's Next method is amoritized O(1), unless the map changes in which case the
// following Next is O(log(n)) where n is the number of elements in the map.
func (c *SetCursor[O, T]) Backward() iterator.Iterator[T] {
	return iterator.Map(c.inner.Backward(), func(kv KVPair[T, struct{}]) T {
		return kv.K
	})
}
