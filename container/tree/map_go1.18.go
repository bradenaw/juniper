//go:build go1.18

package tree

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

type KVPair[K any, V any] struct {
	Key   K
	Value V
}

// Map is a tree-structured key-value map, similar to Go's built-in map but keeps elements in sorted
// order by key.
type Map[K any, V any] struct {
	// An extra indirect here so that tree.Map behaves like a reference type like the map builtin.
	t *btree[K, V]
}

// NewMap returns a Map that uses less to determine the sort order of keys. If !less(a, b) &&
// !less(b, a), then a and b are considered the same key. The output of less must not change for any
// pair of keys while they are in the map.
func NewMap[K any, V any](less xsort.Less[K]) Map[K, V] {
	return Map[K, V]{
		t: newBtree[K, V](less),
	}
}

// Len returns the number of elements in the map.
func (m Map[K, V]) Len() int {
	return m.t.size
}

// Put inserts the key-value pair into the map, overwriting the value for the key if it already
// exists.
func (m Map[K, V]) Put(k K, v V) {
	m.t.Put(k, v)
}

// Delete removes the given key from the map.
func (m Map[K, V]) Delete(k K) {
	m.t.Delete(k)
}

// Get returns the value associated with the given key if it is present in the map. Otherwise, it
// returns the zero-value of V.
func (m Map[K, V]) Get(k K) V {
	return m.t.Get(k)
}

// Contains returns true if the given key is present in the map.
func (m Map[K, V]) Contains(k K) bool {
	return m.t.Contains(k)
}

// First returns the lowest-keyed entry in the map according to less.
func (m Map[K, V]) First() (K, V) {
	return m.t.First()
}

// Last returns the highest-keyed entry in the map according to less.
func (m Map[K, V]) Last() (K, V) {
	return m.t.Last()
}

// Iterate returns an iterator that yields the elements of the map in sorted order by key.
//
// The map may be safely modified during iteration and the iterator will continue from the
// next-lowest key. Thus if the map is modified, the iterator will not necessarily return all of
// the keys present in the map.
func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return m.Cursor().Forward()
}

func (m Map[K, V]) Range(
	lower xsort.Bound[K],
	upper xsort.Bound[K],
	direction xsort.Direction,
) iterator.Iterator[KVPair[K, V]] {
	c := m.Cursor()

	switch direction {
	case xsort.Asc:
		switch lower := lower.(type) {
		case xsort.Min[K]:
		case xsort.Before[K]:
			c.SeekFirstGreaterOrEqual(lower.T)
		case xsort.After[K]:
			c.SeekFirstGreater(lower.T)
		case xsort.Max[K]:
			c.SeekLast()
		default:
			panic("")
		}

		var while func(pair KVPair[K, V]) bool
		switch upper := upper.(type) {
		case xsort.Min[K]:
			while = func(pair KVPair[K, V]) bool { return false }
		case xsort.Before[K]:
			while = func(pair KVPair[K, V]) bool {
				return m.t.less(pair.Key, upper.T)
			}
		case xsort.After[K]:
			while = func(pair KVPair[K, V]) bool {
				return xsort.LessOrEqual(m.t.less, pair.Key, upper.T)
			}
		case xsort.Max[K]:
			while = func(pair KVPair[K, V]) bool { return true }
		default:
			panic("")
		}

		return iterator.While(c.Forward(), while)
	case xsort.Desc:
		switch upper := upper.(type) {
		case xsort.Min[K]:
		case xsort.Before[K]:
			c.SeekLastLess(upper.T)
		case xsort.After[K]:
			c.SeekLastLessOrEqual(upper.T)
		case xsort.Max[K]:
			c.SeekLast()
		default:
			panic("")
		}

		var while func(pair KVPair[K, V]) bool
		switch lower := lower.(type) {
		case xsort.Min[K]:
			while = func(pair KVPair[K, V]) bool { return true }
		case xsort.Before[K]:
			while = func(pair KVPair[K, V]) bool {
				return xsort.GreaterOrEqual(m.t.less, pair.Key, lower.T)
			}
		case xsort.After[K]:
			while = func(pair KVPair[K, V]) bool {
				return xsort.Greater(m.t.less, pair.Key, lower.T)
			}
		case xsort.Max[K]:
			while = func(pair KVPair[K, V]) bool { return false }
		default:
			panic("")
		}

		return iterator.While(c.Backward(), while)
	default:
		panic("")
	}
}

// Cursor returns a cursor into the map placed at the first element.
func (m Map[K, V]) Cursor() *MapCursor[K, V] {
	inner := m.t.Cursor()
	return &MapCursor[K, V]{
		inner: inner,
	}
}

// MapCursor is a cursor into a Map.
//
// A cursor is usable while a map is being modified. If the element the cursor is at is deleted, the
// cursor will still return the old value.
type MapCursor[K any, V any] struct {
	inner cursor[K, V]
}

// SeekFirst moves the cursor to the first element in the map.
//
// SeekFirst is O(log(n)).
func (c *MapCursor[K, V]) SeekFirst() { c.inner.SeekFirst() }

// SeekLast moves the cursor to the last element in the map.
//
// SeekLast is O(log(n)).
func (c *MapCursor[K, V]) SeekLast() { c.inner.SeekLast() }

// SeekLastLess moves the cursor to the element in the map just before k.
//
// SeekLastLess is O(log(n)).
func (c *MapCursor[K, V]) SeekLastLess(k K) { c.inner.SeekLastLess(k) }

// SeekLastLessOrEqual moves the cursor to the element in the map with the greatest key that is less
// than or equal to k.
//
// SeekLastLessOrEqual is O(log(n)).
func (c *MapCursor[K, V]) SeekLastLessOrEqual(k K) { c.inner.SeekLastLessOrEqual(k) }

// SeekFirstGreaterOrEqual moves the cursor to the element in the map with the least key that is
// greater than or equal to k.
//
// SeetFirstGreaterOrEqual is O(log(n)).
func (c *MapCursor[K, V]) SeekFirstGreaterOrEqual(k K) { c.inner.SeekFirstGreaterOrEqual(k) }

// SeekFirstGreater moves the cursor to the element in the map just after k.
//
// SeekFirstGreater is O(log(n)).
func (c *MapCursor[K, V]) SeekFirstGreater(k K) { c.inner.SeekFirstGreater(k) }

// Next moves the cursor to the next element in the map.
//
// Next is amoritized O(1) unless the map has been modified since the last cursor move, in which
// case it's O(log(n)).
func (c *MapCursor[K, V]) Next() { c.inner.Next() }

// Prev moves the cursor to the previous element in the map.
//
// Prev is amoritized O(1) unless the map has been modified since the last cursor move, in which
// case it's O(log(n)).
func (c *MapCursor[K, V]) Prev() { c.inner.Prev() }

// Ok returns false if the cursor is not currently placed at an element, for example if Next
// advances past the last element.
func (c *MapCursor[K, V]) Ok() bool { return c.inner.Ok() }

// Key returns the key of the element that the cursor is at. Panics if Ok is false.
func (c *MapCursor[K, V]) Key() K { return c.inner.Key() }

// Value returns the value of the element that the cursor is at. Panics if Ok is false.
func (c *MapCursor[K, V]) Value() V { return c.inner.Value() }

// Forward returns an iterator that starts from the cursor's position and yields all of the elements
// greater than or equal to the cursor in ascending order.
//
// This iterator's Next method is amoritized O(1), unless the map changes in which case the
// following Next is O(log(n)) where n is the number of elements in the map.
func (c *MapCursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]] { return c.inner.Forward() }

// Backward returns an iterator that starts from the cursor's position and yields all of the
// elements less than or equal to the cursor in descending order.
//
// This iterator's Next method is amoritized O(1), unless the map changes in which case the
// following Next is O(log(n)) where n is the number of elements in the map.
func (c *MapCursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]] { return c.inner.Backward() }
