//go:build go1.18

package tree

import (
	"github.com/bradenaw/xstd/iterator"
	"github.com/bradenaw/xstd/xsort"
)

type KVPair[K any, V any] struct {
	K K
	V V
}

// Map is a tree-structured key-value map, similar to Go's built-in map but keeps elements in sorted
// order by key.
type Map[K any, V any] struct {
	t *tree[KVPair[K, V]]
}

// NewMap returns a Map that uses less to determine the sort order of keys. If !less(a, b) &&
// !less(b, a), then a and b are considered the same key.
func NewMap[K any, V any](less xsort.Less[K]) Map[K, V] {
	return Map[K, V]{
		t: newTree(func(a, b KVPair[K, V]) bool {
			return less(a.K, b.K)
		}),
	}
}

// Len returns the number of elements in the map.
func (m Map[K, V]) Len() int {
	return m.t.size
}

// Put inserts the key-value pair into the map, overwriting the value for the key if it already
// exists.
func (m Map[K, V]) Put(k K, v V) {
	m.t.Put(KVPair[K, V]{k, v})
}

// Delete removes the given key from the map.
func (m Map[K, V]) Delete(k K) {
	m.t.Delete(KVPair[K, V]{K: k})
}

// Contains returns true if the given key is present in the map.
func (m Map[K, V]) Contains(k K) bool {
	return m.t.Contains(KVPair[K, V]{K: k})
}

// Iterate returns an iterator that yields the elements of the map in sorted order by key.
//
// The map may be safely modified during iteration and the iterator will continue from the
// next-lowest key. Thus if the map is modified, the iterator will not necessarily return all of
// the keys present in the map.
func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return m.t.Iterate()
}
