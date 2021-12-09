//go:build go1.18

package tree

import (
	"github.com/bradenaw/xstd/iterator"
	"github.com/bradenaw/xstd/xsort"
)

type Set[T any] struct {
	t *tree[T]
}

func NewSet[T any](less xsort.Less[T]) Set[T] {
	return Set[T]{
		t: newTree(less),
	}
}

func (s Set[T]) Len() int {
	return s.t.size
}

func (s Set[T]) Add(item T) {
	s.t.Put(item)
}

func (s Set[T]) Remove(item T) {
	s.t.Delete(item)
}

func (s Set[T]) Contains(item T) bool {
	return s.t.Contains(item)
}

func (s Set[T]) Iterate() iterator.Iterator[T] {
	return s.t.Iterate()
}
