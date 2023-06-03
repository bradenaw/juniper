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
//
// It is safe for multiple goroutines to Put concurrently with keys that are already in the map.
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

// Iterate returns an iterator that yields the elements of the map in ascending order by key.
//
// The map may be safely modified during iteration and the iterator will continue from the
// next-lowest key. Thus the iterator will see new elements that are after the current position of
// the iterator according to less, but will not necessarily see a consistent snapshot of the state
// of the map.
func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return m.Range(Unbounded[K](), Unbounded[K]())
}

type boundType int

const (
	boundInclude boundType = iota + 1
	boundExclude
	boundUnbounded
)

// Bound is an endpoint for a range.
type Bound[K any] struct {
	type_ boundType
	key   K
}

// Included returns a Bound that goes up to and including key.
func Included[K any](key K) Bound[K] { return Bound[K]{type_: boundInclude, key: key} }

// Excluded returns a Bound that goes up to but not including key.
func Excluded[K any](key K) Bound[K] { return Bound[K]{type_: boundExclude, key: key} }

// Unbounded returns a Bound at the end of the collection.
func Unbounded[K any]() Bound[K] { return Bound[K]{type_: boundUnbounded} }

// Range returns an iterator that yields the elements of the map between the given bounds in
// ascending order by key.
//
// The map may be safely modified during iteration and the iterator will continue from the
// next-lowest key. Thus the iterator will see new elements that are after the current position of
// the iterator according to less, but will not necessarily see a consistent snapshot of the state
// of the map.
func (m Map[K, V]) Range(lower Bound[K], upper Bound[K]) iterator.Iterator[KVPair[K, V]] {
	return m.t.Range(lower, upper)
}

// RangeReverse returns an iterator that yields the elements of the map between the given bounds in
// descending order by key.
//
// The map may be safely modified during iteration and the iterator will continue from the
// next-lowest key. Thus the iterator will see new elements that are after the current position of
// the iterator according to less, but will not necessarily see a consistent snapshot of the state
// of the map.
func (m Map[K, V]) RangeReverse(lower Bound[K], upper Bound[K]) iterator.Iterator[KVPair[K, V]] {
	return m.t.RangeReverse(lower, upper)
}
