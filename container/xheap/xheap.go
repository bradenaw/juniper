package xheap

import (
	"fmt"

	"github.com/bradenaw/xstd/internal/heap"
	"github.com/bradenaw/xstd/xsort"
)

type Heap[T any] struct {
	// Indirect here so that Heap behaves as a reference type, like the map builtin.
	inner *heap.Heap[T]
}

func New[T any](less xsort.Less[T], initial []T) Heap[T] {
	inner := heap.New(
		func(a, b T) bool {
			return less(a, b)
		},
		func(a T, i int) {},
		initial,
	)
	return Heap[T]{
		inner: &inner,
	}
}

func (h *Heap[T]) Len() int {
	return h.inner.Len()
}

func (h *Heap[T]) Grow(n int) {
	h.inner.Grow(n)
}

func (h *Heap[T]) Push(item T) {
	h.inner.Push(item)
}

func (h *Heap[T]) Pop() (T, bool) {
	return h.inner.Pop()
}

type KVPair[K any, V any] struct {
	K K
	V V
}

type MapHeap[K comparable, V any] struct {
	// Indirect here so that Heap behaves as a reference type, like the map builtin.
	inner *heap.Heap[KVPair[K, V]]
	m     map[K]int
}

func NewMap[K comparable, V any](less xsort.Less[K], initial []KVPair[K, V]) MapHeap[K, V] {
	h := MapHeap[K, V]{
		m: make(map[K]int),
	}
	inner := heap.New(
		func(a, b KVPair[K, V]) bool {
			return less(a.K, b.K)
		},
		func(x KVPair[K, V], i int) {
			h.m[x.K] = i
		},
		initial,
	)
	h.inner = &inner
	return h
}

func (h *MapHeap[K, V]) Len() int {
	return h.inner.Len()
}

func (h *MapHeap[K, V]) Grow(n int) {
	h.inner.Grow(n)
}

func (h *MapHeap[K, V]) Push(k K, v V) {
	h.inner.Push(KVPair[K, V]{k, v})
}

func (h *MapHeap[K, V]) Pop() (KVPair[K, V], bool) {
	item, ok := h.inner.Pop()
	if ok {
		delete(h.m, item.K)
	}
	return item, ok
}

func (h *MapHeap[K, V]) Remove(k K) {
	i, ok := h.m[k]
	if !ok {
		panic(fmt.Sprintf("remove item not in MapHeap: %#v", k))
	}
	h.inner.RemoveAt(i)
	delete(h.m, k)
}
