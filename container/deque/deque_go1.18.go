//go:build go1.18

// Package deque contains a double-ended queue.
package deque

import (
	"errors"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xmath"
)

var errDequeEmpty = errors.New("pop from empty deque")
var errDequeModified = errors.New("deque modified during iteration")

const (
	// A non-empty deque has space for at least this many items.
	minSize = 16
	// When growing a full deque, reallocate with len(r.a)*growFactor.
	growFactor = 2
	// Shrink a deque that is 1/shrinkFactor full, down to 1/shrinkFactor size.
	shrinkFactor = 16
)

// Deque is a double-ended queue, allowing push and pop to both the front and back of the queue.
// Pushes and pops are amoritized O(1). The zero-value is ready to use. Deque should not be copied
// after first use.
type Deque[T any] struct {
	// Backing slice for the deque. Empty if the deque is empty.
	a []T
	// Index of the first item.
	front int
	// Index of the last item.
	back int
	gen  int
}

// Len returns the number of items in the deque.
func (r *Deque[T]) Len() int {
	if len(r.a) == 0 {
		return 0
	}

	if r.front <= r.back {
		return r.back - r.front + 1
	}
	return len(r.a) - r.front + r.back + 1
}

// Grow allocates sufficient space to add n more items without needing to reallocate.
func (r *Deque[T]) Grow(n int) {
	extraCap := len(r.a) - r.Len()
	if extraCap < n {
		r.resize(len(r.a) + n)
	}
}

// PushFront adds item to the front of the deque.
func (r *Deque[T]) PushFront(item T) {
	if len(r.a) == 0 {
		r.a = make([]T, minSize)
		r.a[0] = item
		return
	}
	r.maybeExpand()
	r.front = positiveMod(r.front-1, len(r.a))
	r.a[r.front] = item
	r.gen++
}

// PushFront adds item to the back of the deque.
func (r *Deque[T]) PushBack(item T) {
	if len(r.a) == 0 {
		r.a = make([]T, minSize)
		r.a[0] = item
		return
	}
	r.maybeExpand()
	r.back = (r.back + 1) % len(r.a)
	r.a[r.back] = item
	r.gen++
}

// Guarantees that there is room in the deque.
func (r *Deque[T]) maybeExpand() {
	if r.Len() == len(r.a) {
		r.resize(xmath.Max(minSize, len(r.a)*2))
	}
}

func (r *Deque[T]) maybeShrink() {
	l := r.Len()
	if l > minSize && l < len(r.a)/shrinkFactor {
		r.resize(len(r.a) / shrinkFactor)
	}
}

func (r *Deque[T]) resize(n int) {
	oldLen := r.Len()
	newA := make([]T, n)
	copy(newA[0:], r.a[r.front:])
	copy(newA[len(r.a)-r.front:], r.a[:r.front])
	r.a = newA
	r.front = 0
	r.back = oldLen - 1
}

// PopFront removes and returns the item at the front of the deque. It panics if the deque is empty.
func (r *Deque[T]) PopFront() T {
	l := r.Len()
	if l == 0 {
		panic(errDequeEmpty)
	}
	item := r.a[r.front]
	if l == 1 {
		r.a = nil
		r.front = 0
		r.back = 0
		return item
	}
	var zero T
	r.a[r.front] = zero
	r.front = (r.front + 1) % len(r.a)
	r.maybeShrink()
	r.gen++
	return item
}

// PopBack removes and returns the item at the back of the deque. It panics if the deque is empty.
func (r *Deque[T]) PopBack() T {
	l := r.Len()
	if l == 0 {
		panic(errDequeEmpty)
	}
	item := r.a[r.back]
	if l == 1 {
		r.a = nil
		r.front = 0
		r.back = 0
		return item
	}
	var zero T
	r.a[r.back] = zero
	r.back = positiveMod(r.back-1, len(r.a))
	r.maybeShrink()
	r.gen++
	return item
}

// PeekFront returns the item at the front of the deque. It panics if the deque is empty.
func (r *Deque[T]) PeekFront() T {
	return r.a[r.front]
}

// PeekBack returns the item at the back of the deque. It panics if the deque is empty.
func (r *Deque[T]) PeekBack() T {
	return r.a[r.back]
}

func positiveMod(l, r int) int {
	x := l % r
	if x < 0 {
		return x + r
	}
	return x
}

// Iterate iterates over the elements of the deque.
//
// The iterator panics if the deque has been modified since iteration started.
func (r *Deque[T]) Iterate() iterator.Iterator[T] {
	i := r.front
	done := false
	gen := r.gen
	return iterator.FromNext(func() (T, bool) {
		if gen != r.gen {
			panic(errDequeModified)
		}
		var zero T
		if r.Len() == 0 {
			return zero, false
		}
		if done {
			return zero, false
		}
		item := r.a[i]
		if i == r.back {
			done = true
		}
		i = (i + 1) % len(r.a)
		return item, true
	})
}
