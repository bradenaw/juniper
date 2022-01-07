//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

type KVPair[K any, V any] struct {
	K K
	V V
}

type KeyOrder[K any, V any, O xsort.Ordering[K]] struct{}

func (KeyOrder[K, V, O]) Less(a, b KVPair[K, V]) bool {
	var ordering O
	return ordering.Less(a.K, b.K)
}

// Map is a tree-structured key-value map, similar to Go's built-in map but keeps elements in sorted
// order by key.
type Map[K any, V any, O xsort.Ordering[K]] struct {
	t *tree[K, V, O]
}

// NewMap returns a Map that uses less to determine the sort order of keys. If !less(a, b) &&
// !less(b, a), then a and b are considered the same key.
func NewMap[K any, V any, O xsort.Ordering[K]]() Map[K, V, O] {
	return Map[K, V, O]{
		t: newTree[K, V, O](),
	}
}

// Len returns the number of elements in the map.
func (m Map[K, V, O]) Len() int {
	return m.t.size
}

// Put inserts the key-value pair into the map, overwriting the value for the key if it already
// exists.
func (m Map[K, V, O]) Put(k K, v V) {
	m.t.Put(k, v)
}

// Delete removes the given key from the map.
func (m Map[K, V, O]) Delete(k K) {
	m.t.Delete(k)
}

// Get returns the value associated with the given key if it is present in the map. Otherwise, it
// returns the zero-value of V.
func (m Map[K, V, O]) Get(k K) V {
	return m.t.Get(k)
}

// Contains returns true if the given key is present in the map.
func (m Map[K, V, O]) Contains(k K) bool {
	return m.t.Contains(k)
}

// First returns the lowest-keyed entry in the map according to less.
func (m Map[K, V, O]) First() (K, V) {
	return m.t.First()
}

// Last returns the highest-keyed entry in the map according to less.
func (m Map[K, V, O]) Last() (K, V) {
	return m.t.Last()
}

// Iterate returns an iterator that yields the elements of the map in sorted order by key.
//
// The map may be safely modified during iteration and the iterator will continue from the
// next-lowest key. Thus if the map is modified, the iterator will not necessarily return all of
// the keys present in the map.
func (m Map[K, V, O]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return m.Cursor().Forward()
}

// Cursor returns a cursor into the map placed at the first element.
func (m Map[K, V, O]) Cursor() *MapCursor[K, V, O] {
	inner := m.t.Cursor()
	return &MapCursor[K, V, O]{
		inner: inner,
	}
}

// MapCursor is a cursor into a Map.
//
// A cursor is usable while a map is being modified. If the element the cursor is at is deleted, the
// cursor will still return the old value.
type MapCursor[K any, V any, O xsort.Ordering[K]] struct {
	inner cursor[K, V, O]
}

// SeekFirst moves the cursor to the first element in the map.
//
// SeekFirst is O(log(n)).
func (c *MapCursor[K, V, O]) SeekFirst() { c.inner.SeekFirst() }

// SeekLast moves the cursor to the last element in the map.
//
// SeekLast is O(log(n)).
func (c *MapCursor[K, V, O]) SeekLast() { c.inner.SeekLast() }

// SeekLastLess moves the cursor to the element in the map just before k.
//
// SeekLastLess is O(log(n)).
func (c *MapCursor[K, V, O]) SeekLastLess(k K) { c.inner.SeekLastLess(k) }

// SeekLastLessOrEqual moves the cursor to the element in the map with the greatest key that is less
// than or equal to k.
//
// SeekLastLessOrEqual is O(log(n)).
func (c *MapCursor[K, V, O]) SeekLastLessOrEqual(k K) { c.inner.SeekLastLessOrEqual(k) }

// SeekFirstGreaterOrEqual moves the cursor to the element in the map with the least key that is
// greater than or equal to k.
//
// SeetFirstGreaterOrEqual is O(log(n)).
func (c *MapCursor[K, V, O]) SeekFirstGreaterOrEqual(k K) { c.inner.SeekFirstGreaterOrEqual(k) }

// SeekFirstGreater moves the cursor to the element in the map just after k.
//
// SeekFirstGreater is O(log(n)).
func (c *MapCursor[K, V, O]) SeekFirstGreater(k K) { c.inner.SeekFirstGreater(k) }

// Next moves the cursor to the next element in the map.
//
// Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
// case it's O(log(n)).
func (c *MapCursor[K, V, O]) Next() { c.inner.Next() }

// Prev moves the cursor to the previous element in the map.
//
// Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
// case it's O(log(n)).
func (c *MapCursor[K, V, O]) Prev() { c.inner.Prev() }

// Ok returns false if the cursor is not currently placed at an element, for example if Next
// advances past the last element.
func (c *MapCursor[K, V, O]) Ok() bool { return c.inner.Ok() }

// Key returns the key of the element that the cursor is at. Panics if Ok is false.
func (c *MapCursor[K, V, O]) Key() K { return c.inner.Key() }

// Value returns the value of the element that the cursor is at. Panics if Ok is false.
func (c *MapCursor[K, V, O]) Value() V { return c.inner.Value() }

// Forward returns an iterator that starts from the cursor's position and yields all of the elements
// greater than or equal to the cursor in ascending order.
//
// This iterator's Next method is amoritized O(1), unless the map changes in which case the
// following Next is O(log(n)) where n is the number of elements in the map.
func (c *MapCursor[K, V, O]) Forward() iterator.Iterator[KVPair[K, V]] { return c.inner.Forward() }

// Backward returns an iterator that starts from the cursor's position and yields all of the
// elements less than or equal to the cursor in descending order.
//
// This iterator's Next method is amoritized O(1), unless the map changes in which case the
// following Next is O(log(n)) where n is the number of elements in the map.
func (c *MapCursor[K, V, O]) Backward() iterator.Iterator[KVPair[K, V]] { return c.inner.Backward() }
