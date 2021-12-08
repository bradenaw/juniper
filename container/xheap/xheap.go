package xheap

import (
	"github.com/bradenaw/xstd/internal/heap"
	"github.com/bradenaw/xstd/xsort"
)

type Heap[T any] struct {
	// Indirect here so that Heap behaves as a reference type, like the map builtin.
	*heap.Heap[T]
}

func New[T any](less xsort.Less[T], initial []T) Heap[T] {
	inner := heap.New(func(a, b T) bool {
		return less(a, b)
	}, initial)
	return Heap[T]{
		Heap: &inner,
	}
}
