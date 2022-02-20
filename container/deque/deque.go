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
// Pushes and pops are amortized O(1). The zero-value is ready to use. Deque should not be copied
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
	if r.a == nil || r.back == -1 {
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
	if r.Len() == 0 {
		r.a = make([]T, minSize)
		r.a[0] = item
		r.back = 0
		return
	}
	r.maybeExpand()
	r.front = positiveMod(r.front-1, len(r.a))
	r.a[r.front] = item
	r.gen++
}

// PushFront adds item to the back of the deque.
func (r *Deque[T]) PushBack(item T) {
	if r.Len() == 0 {
		r.a = make([]T, minSize)
		r.a[0] = item
		r.back = 0
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
	if !(r.a == nil || r.back == -1) {
		if r.front <= r.back {
			copy(newA, r.a[r.front:r.back+1])
		} else {
			copy(newA, r.a[r.front:])
			copy(newA[len(r.a)-r.front:], r.a[:r.back+1])
		}
	}
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
		r.back = -1
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
		r.back = -1
		return item
	}
	var zero T
	r.a[r.back] = zero
	r.back = positiveMod(r.back-1, len(r.a))
	r.maybeShrink()
	r.gen++
	return item
}

// Front returns the item at the front of the deque. It panics if the deque is empty.
func (r *Deque[T]) Front() T {
	if r.back == -1 {
		panic("deque index out of range")
	}
	return r.a[r.front]
}

// Back returns the item at the back of the deque. It panics if the deque is empty.
func (r *Deque[T]) Back() T {
	return r.a[r.back]
}

// Item returns the ith item in the deque. 0 is the front and r.Len()-1 is the back.
func (r *Deque[T]) Item(i int) T {
	if i < 0 || i >= r.Len() {
		panic("deque index out of range")
	}
	idx := (r.front + i) % len(r.a)
	return r.a[idx]
}

func positiveMod(l, r int) int {
	x := l % r
	if x < 0 {
		return x + r
	}
	return x
}

type dequeIterator[T any] struct {
	r    *Deque[T]
	i    int
	done bool
	gen  int
}

func (iter *dequeIterator[T]) Next() (T, bool) {
	if iter.gen != iter.r.gen {
		panic(errDequeModified)
	}
	var zero T
	if iter.r.Len() == 0 {
		return zero, false
	}
	if iter.done {
		return zero, false
	}
	item := iter.r.a[iter.i]
	if iter.i == iter.r.back {
		iter.done = true
	}
	iter.i = (iter.i + 1) % len(iter.r.a)
	return item, true
}

// Iterate iterates over the elements of the deque.
//
// The iterator panics if the deque has been modified since iteration started.
func (r *Deque[T]) Iterate() iterator.Iterator[T] {
	return &dequeIterator[T]{
		r:    r,
		i:    r.front,
		done: false,
		gen:  r.gen,
	}
}
