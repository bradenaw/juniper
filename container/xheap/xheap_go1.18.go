//go:build go1.18

// Package xheap contains extensions to the standard library package container/heap.
package xheap

import (
	"fmt"

	"github.com/bradenaw/juniper/internal/heap"
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

// Heap is a min-heap (https://en.wikipedia.org/wiki/Binary_heap). Min-heaps are a collection
// structure that provide constant-time access to the minimum element, and logarithmic-time removal.
// They are most commonly used as a priority queue.
type Heap[T any] struct {
	// Indirect here so that Heap behaves as a reference type, like the map builtin.
	inner *heap.Heap[T]
}

// New returns a new Heap which uses less to determine the minimum element.
//
// The elements from initial are added to the heap. initial is modified by New and utilized by the
// Heap, so it should not be used after passing to New(). Passing initial is faster (O(n)) than
// creating an empty heap and pushing each item (O(n * log(n))).
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

// Len returns the current number of elements in the heap.
func (h *Heap[T]) Len() int {
	return h.inner.Len()
}

// Grow allocates sufficient space to add n more elements without needing to reallocate.
func (h *Heap[T]) Grow(n int) {
	h.inner.Grow(n)
}

// Push adds item to the heap.
func (h *Heap[T]) Push(item T) {
	h.inner.Push(item)
}

// Pop removes and returns the minimum item in the heap. It panics if h.Len()==0.
func (h *Heap[T]) Pop() T {
	return h.inner.Pop()
}

// Pop returns the minimum item in the heap. It panics if h.Len()==0.
func (h *Heap[T]) Peek() T {
	return h.inner.Peek()
}

// Iterate iterates over the elements of the heap.
//
// The iterator panics if the heap has been modified since iteration started.
func (h *Heap[T]) Iterate() iterator.Iterator[T] {
	return h.inner.Iterate()
}

type KVPair[K any, V any] struct {
	K K
	V V
}

// MapHeap is equivalent to Heap except that it stores key-value pairs. Only the key is used to
// determine the minimum element. Only one element with the same key can be in the MapHeap at a
// time.
type MapHeap[K comparable, V any] struct {
	// Indirect here so that Heap behaves as a reference type, like the map builtin.
	inner *heap.Heap[KVPair[K, V]]
	m     map[K]int
}

// NewMap returns a new MapHeap which uses less to determine the minimum element.
//
// The elements from initial are added to the heap. initial is modified by New and utilized by the
// Heap, so it should not be used after passing to New(). Passing initial is faster (O(n)) than
// creating an empty heap and pushing each item (O(n * log(n))).
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

// Len returns the current number of elements in the heap.
func (h *MapHeap[K, V]) Len() int {
	return h.inner.Len()
}

// Grow allocates sufficient space to add n more elements without needing to reallocate.
func (h *MapHeap[K, V]) Grow(n int) {
	h.inner.Grow(n)
}

// Push adds item to the heap.
func (h *MapHeap[K, V]) Push(k K, v V) {
	h.inner.Push(KVPair[K, V]{k, v})
}

// Pop removes and returns the minimum item in the heap. It panics if h.Len()==0.
func (h *MapHeap[K, V]) Pop() KVPair[K, V] {
	item := h.inner.Pop()
	delete(h.m, item.K)
	return item
}

// Pop returns the minimum item in the heap. It panics if h.Len()==0.
func (h *MapHeap[K, V]) Peek() KVPair[K, V] {
	return h.inner.Peek()
}

// Contains returns true if the given key is present in the heap.
func (h *MapHeap[K, V]) Contains(k K) bool {
	_, ok := h.m[k]
	return ok
}

// Remove removes the item with the given key. It panics if the key is not present in the heap.
func (h *MapHeap[K, V]) Remove(k K) {
	i, ok := h.m[k]
	if !ok {
		panic(fmt.Sprintf("remove item not in MapHeap: %#v", k))
	}
	h.inner.RemoveAt(i)
	delete(h.m, k)
}

// Iterate iterates over the elements of the heap.
//
// The iterator panics if the heap has been modified since iteration started.
func (h *MapHeap[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return h.inner.Iterate()
}
