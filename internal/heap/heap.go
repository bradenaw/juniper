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

type Heap[T any, O Ordering[T]] struct {
	indexChanged func(x T, i int)
	a            []T
	gen          int
}

func New[T any, O Ordering[T]](indexChanged func(x T, i int), initial []T) Heap[T, O] {
	h := Heap[T, O]{
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

func (h *Heap[T, O]) Len() int {
	return len(h.a)
}

func (h *Heap[T, O]) Grow(n int) {
	h.a = slices.Grow(h.a, n)
}

func (h *Heap[T, O]) Push(item T) {
	h.a = append(h.a, item)
	h.notifyIndexChanged(len(h.a) - 1)
	h.percolateUp(len(h.a) - 1)
	h.gen++
}

func (h *Heap[T, O]) Pop() T {
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

func (h *Heap[T, O]) Peek() T {
	return h.a[0]
}

func (h *Heap[T, O]) RemoveAt(i int) {
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

func (h *Heap[T, O]) Item(i int) T {
	return h.a[i]
}

func (h *Heap[T, O]) UpdateAt(i int, item T) {
	h.a[i] = item
	h.notifyIndexChanged(i)
	h.percolateUp(i)
	h.percolateDown(i)
}

func (h *Heap[T, O]) maybeShrink() {
	if len(h.a) > 0 && cap(h.a)/len(h.a) >= shrinkFactor {
		newA := make([]T, len(h.a))
		copy(newA, h.a)
		h.a = newA
	}
}

func (h *Heap[T, O]) percolateUp(i int) {
	for i > 0 {
		p := parent(i)
		if h.lessByIdx(i, p) {
			h.swap(i, p)
		}
		i = p
	}
}

func (h *Heap[T, O]) swap(i, j int) {
	(h.a)[i], (h.a)[j] = (h.a)[j], (h.a)[i]
	h.notifyIndexChanged(i)
	h.notifyIndexChanged(j)
}

func (h *Heap[T, O]) notifyIndexChanged(i int) {
	h.indexChanged(h.a[i], i)
}

func (h *Heap[T, O]) lessByIdx(i, j int) bool {
	var ordering O
	return ordering.Less((h.a)[i], (h.a)[j])
}

func (h *Heap[T, O]) percolateDown(i int) {
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

type heapIterator[T any, O Ordering[T]] struct {
	h     *Heap[T, O]
	inner iterator.Iterator[T]
	gen   int
}

func (iter *heapIterator[T, O]) Next() (T, bool) {
	if iter.gen == -1 {
		iter.gen = iter.h.gen
		iter.inner = iterator.Slice(iter.h.a)
	} else if iter.gen != iter.h.gen {
		panic(ErrHeapModified)
	}
	return iter.inner.Next()
}

func (h *Heap[T, O]) Iterate() iterator.Iterator[T] {
	return &heapIterator[T, O]{h: h, gen: -1}
}

func parent(i int) int {
	return (i - 1) / 2
}

func children(i int) (int, int) {
	return i*2 + 1, i*2 + 2
}
