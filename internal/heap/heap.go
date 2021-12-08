package heap

import (
	"github.com/bradenaw/xstd/iterator"
	"github.com/bradenaw/xstd/slices"
)

// Duplicated from xsort to avoid dependency cycle.
type Less[T any] func(a, b T) bool

type Heap[T any] struct {
	lessFn Less[T]
	// Indirect here to make Heap a reference type like the built-in Map.
	a *[]T
}

func New[T any](less Less[T], initial []T) Heap[T] {
	h := Heap[T]{
		lessFn: less,
		a: &initial,
	}

	for i := len(initial) / 2 - 1; i >= 0; i-- {
		h.percolateDown(i)
	}

	return h
}

func (h *Heap[T]) Len() int {
	return len(*h.a)
}

func (h *Heap[T]) Grow(n int) {
	*h.a = slices.Grow(*h.a, n)
}

func (h *Heap[T]) Push(item T) {
	*h.a = append(*h.a, item)
	h.percolateUp(len(*h.a)-1)
}

func (h *Heap[T]) Pop() (T, bool) {
	var zero T
	if len(*h.a) == 0 {
		return zero, false
	}
	item := (*h.a)[0]
	(*h.a)[0] = (*h.a)[len(*h.a)-1]
	// In case T is a pointer, clear this out to keep the ref from being live.
	(*h.a)[len(*h.a)-1] = zero
	*h.a = (*h.a)[:len(*h.a)-1]
	h.percolateDown(0)
	return item, true
}

func (h *Heap[T]) percolateUp(i int) {
	for i > 0 {
		p := parent(i)
		if h.less(i, p) {
			h.swap(i, p)
		}
		i = p
	}
}

func (h *Heap[T]) swap(i, j int) {
	(*h.a)[i], (*h.a)[j] = (*h.a)[j], (*h.a)[i]
}

func (h *Heap[T]) less(i, j int) bool {
	return h.lessFn((*h.a)[i], (*h.a)[j])
}

func (h *Heap[T]) percolateDown(i int) {
	for {
		left, right := children(i)
		if left >= len(*h.a) {
			// no children
			return
		} else if right >= len(*h.a) {
			// only has a left child
			if h.less(left, i) {
				h.swap(left, i)
				i = left
			} else {
				return
			}
		} else {
			// has both children
			least := left
			if h.less(right, left) {
				least = right
			}
			if h.less(least, i) {
				h.swap(least, i)
				i = least
			} else {
				return
			}
		}
	}
}

func (h *Heap[T]) Iterate() iterator.Iterator[T] {
	return iterator.Slice(*h.a)
}

func parent(i int) int {
	return (i - 1) / 2
}

func children(i int) (int, int) {
	return i * 2 + 1, i * 2 + 2
}
