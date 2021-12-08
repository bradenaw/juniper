package tree

import (
	"github.com/bradenaw/xstd/iterator"
	"github.com/bradenaw/xstd/xsort"
)

type KVPair[K any, V any] struct {
	k K
	v V
}

type Map[K any, V any] struct {
	t *tree[KVPair[K, V]]
}

func NewMap[K any, V any](less xsort.Less[K]) Map[K, V] {
	return Map[K, V]{
		t: newTree(func(a, b KVPair[K, V]) bool {
			return less(a.k, b.k)
		}),
	}
}

func (m Map[K, V]) Len() int {
	return m.t.size
}

func (m Map[K, V]) Put(k K, v V) {
	m.t.Put(KVPair[K, V]{k, v})
}

func (m Map[K, V]) Delete(k K) {
	m.t.Delete(KVPair[K, V]{k: k})
}

func (m Map[K, V]) Contains(k K) bool {
	return m.t.Contains(KVPair[K, V]{k: k})
}

func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return m.t.Iterate()
}
