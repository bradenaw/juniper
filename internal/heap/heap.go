//go:build go1.18

package heap

import (
	"errors"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
)

var ErrHeapModified = errors.New("heap modified during iteration")

const shrinkFactor = 16

// Duplicated from xsort to avoid dependency cycle.
type Ordering[T any] interface {
	Less(a, b T) bool
}

type Heap[O Ordering[T], T any] struct {
	indexChanged func(x T, i int)
	a            []T
	gen          int
}

func New[O Ordering[T], T any](indexChanged func(x T, i int), initial []T) Heap[O, T] {
	h := Heap[O, T]{
		indexChanged: indexChanged,
		a:            initial,
	}

	for i := len(initial)/2 - 1; i >= 0; i-- {
		h.percolateDown(i)
	}
	for i := range initial {
		h.notifyIndexChanged(i)
	}

	return h
}

func (h *Heap[O, T]) Len() int {
	return len(h.a)
}

func (h *Heap[O, T]) Grow(n int) {
	h.a = slices.Grow(h.a, n)
}

func (h *Heap[O, T]) Push(item T) {
	h.a = append(h.a, item)
	h.notifyIndexChanged(len(h.a) - 1)
	h.percolateUp(len(h.a) - 1)
	h.gen++
}

func (h *Heap[O, T]) Pop() T {
	var zero T
	item := h.a[0]
	(h.a)[0] = (h.a)[len(h.a)-1]
	// In case T is a pointer, clear this out to keep the ref from being live.
	(h.a)[len(h.a)-1] = zero
	h.a = (h.a)[:len(h.a)-1]
	if len(h.a) > 0 {
		h.notifyIndexChanged(0)
	}
	h.percolateDown(0)
	h.maybeShrink()
	h.gen++
	return item
}

func (h *Heap[O, T]) Peek() T {
	return h.a[0]
}

func (h *Heap[O, T]) RemoveAt(i int) {
	var zero T
	h.a[i] = h.a[len(h.a)-1]
	h.a[len(h.a)-1] = zero
	h.a = h.a[:len(h.a)-1]
	if i < len(h.a) {
		h.notifyIndexChanged(i)
		h.percolateUp(i)
		h.percolateDown(i)
	}
	h.maybeShrink()
	h.gen++
}

func (h *Heap[O, T]) Item(i int) T {
	return h.a[i]
}

func (h *Heap[O, T]) UpdateAt(i int, item T) {
	h.a[i] = item
	h.notifyIndexChanged(i)
	h.percolateUp(i)
	h.percolateDown(i)
}

func (h *Heap[O, T]) maybeShrink() {
	if len(h.a) > 0 && cap(h.a)/len(h.a) >= shrinkFactor {
		newA := make([]T, len(h.a))
		copy(newA, h.a)
		h.a = newA
	}
}

func (h *Heap[O, T]) percolateUp(i int) {
	for i > 0 {
		p := parent(i)
		if h.lessByIdx(i, p) {
			h.swap(i, p)
		}
		i = p
	}
}

func (h *Heap[O, T]) swap(i, j int) {
	(h.a)[i], (h.a)[j] = (h.a)[j], (h.a)[i]
	h.notifyIndexChanged(i)
	h.notifyIndexChanged(j)
}

func (h *Heap[O, T]) notifyIndexChanged(i int) {
	h.indexChanged(h.a[i], i)
}

func (h *Heap[O, T]) lessByIdx(i, j int) bool {
	var ordering O
	return ordering.Less((h.a)[i], (h.a)[j])
}

func (h *Heap[O, T]) percolateDown(i int) {
	for {
		left, right := children(i)
		if left >= len(h.a) {
			// no children
			return
		} else if right >= len(h.a) {
			// only has a left child
			if h.lessByIdx(left, i) {
				h.swap(left, i)
				i = left
			} else {
				return
			}
		} else {
			// has both children
			least := left
			if h.lessByIdx(right, left) {
				least = right
			}
			if h.lessByIdx(least, i) {
				h.swap(least, i)
				i = least
			} else {
				return
			}
		}
	}
}

type heapIterator[O Ordering[T], T any] struct {
	h     *Heap[O, T]
	inner iterator.Iterator[T]
	gen   int
}

func (iter *heapIterator[O, T]) Next() (T, bool) {
	if iter.gen == -1 {
		iter.gen = iter.h.gen
		iter.inner = iterator.Slice(iter.h.a)
	} else if iter.gen != iter.h.gen {
		panic(ErrHeapModified)
	}
	return iter.inner.Next()
}

func (h *Heap[O, T]) Iterate() iterator.Iterator[T] {
	return &heapIterator[O, T]{h: h, gen: -1}
}

func parent(i int) int {
	return (i - 1) / 2
}

func children(i int) (int, int) {
	return i*2 + 1, i*2 + 2
}
