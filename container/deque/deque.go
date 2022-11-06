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
	// When growing a full deque, reallocate with len(d.a)*growFactor.
	growFactor = 2
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
func (d *Deque[T]) Len() int {
	if d.a == nil || d.back == -1 {
		return 0
	}

	if d.front <= d.back {
		return d.back - d.front + 1
	}
	return len(d.a) - d.front + d.back + 1
}

// Grow allocates sufficient space to add n more items without needing to reallocate.
func (d *Deque[T]) Grow(n int) {
	extraCap := len(d.a) - d.Len()
	if extraCap < n {
		d.resize(len(d.a) + n)
	}
}

// Shrink reallocates the backing buffer for d, if necessary, so that it fits only the current size
// plus at most n extra items.
func (d *Deque[T]) Shrink(n int) {
	if n < 0 {
		panic("Shrink() with a negative number of extras")
	}
	if len(d.a)-d.Len() > n {
		d.resize(d.Len() + n)
	}
}

// PushFront adds item to the front of the deque.
func (d *Deque[T]) PushFront(item T) {
	if d.Len() == 0 {
		d.a = make([]T, minSize)
		d.a[0] = item
		d.back = 0
		return
	}
	d.maybeExpand()
	d.front = positiveMod(d.front-1, len(d.a))
	d.a[d.front] = item
	d.gen++
}

// PushFront adds item to the back of the deque.
func (d *Deque[T]) PushBack(item T) {
	if d.Len() == 0 {
		d.a = make([]T, minSize)
		d.a[0] = item
		d.back = 0
		return
	}
	d.maybeExpand()
	d.back = (d.back + 1) % len(d.a)
	d.a[d.back] = item
	d.gen++
}

// Guarantees that there is room in the deque.
func (d *Deque[T]) maybeExpand() {
	if d.Len() == len(d.a) {
		d.resize(xmath.Max(minSize, len(d.a)*2))
	}
}

func (d *Deque[T]) resize(n int) {
	oldLen := d.Len()
	newA := make([]T, n)
	if !(d.a == nil || d.back == -1) {
		if d.front <= d.back {
			copy(newA, d.a[d.front:d.back+1])
		} else {
			copy(newA, d.a[d.front:])
			copy(newA[len(d.a)-d.front:], d.a[:d.back+1])
		}
	}
	d.a = newA
	d.front = 0
	d.back = oldLen - 1
}

// PopFront removes and returns the item at the front of the deque. It panics if the deque is empty.
func (d *Deque[T]) PopFront() T {
	l := d.Len()
	if l == 0 {
		panic(errDequeEmpty)
	}
	item := d.a[d.front]
	if l == 1 {
		d.a = nil
		d.front = 0
		d.back = -1
		return item
	}
	var zero T
	d.a[d.front] = zero
	d.front = (d.front + 1) % len(d.a)
	d.gen++
	return item
}

// PopBack removes and returns the item at the back of the deque. It panics if the deque is empty.
func (d *Deque[T]) PopBack() T {
	l := d.Len()
	if l == 0 {
		panic(errDequeEmpty)
	}
	item := d.a[d.back]
	if l == 1 {
		d.a = nil
		d.front = 0
		d.back = -1
		return item
	}
	var zero T
	d.a[d.back] = zero
	d.back = positiveMod(d.back-1, len(d.a))
	d.gen++
	return item
}

// Front returns the item at the front of the deque. It panics if the deque is empty.
func (d *Deque[T]) Front() T {
	if d.back == -1 {
		panic("deque index out of range")
	}
	return d.a[d.front]
}

// Back returns the item at the back of the deque. It panics if the deque is empty.
func (d *Deque[T]) Back() T {
	return d.a[d.back]
}

// Item returns the ith item in the deque. 0 is the front and d.Len()-1 is the back.
func (d *Deque[T]) Item(i int) T {
	if i < 0 || i >= d.Len() {
		panic("deque index out of range")
	}
	idx := (d.front + i) % len(d.a)
	return d.a[idx]
}

func positiveMod(l, d int) int {
	x := l % d
	if x < 0 {
		return x + d
	}
	return x
}

type dequeIterator[T any] struct {
	d    *Deque[T]
	i    int
	done bool
	gen  int
}

func (iter *dequeIterator[T]) Next() (T, bool) {
	if iter.gen != iter.d.gen {
		panic(errDequeModified)
	}
	var zero T
	if iter.d.Len() == 0 {
		return zero, false
	}
	if iter.done {
		return zero, false
	}
	item := iter.d.a[iter.i]
	if iter.i == iter.d.back {
		iter.done = true
	}
	iter.i = (iter.i + 1) % len(iter.d.a)
	return item, true
}

// Iterate iterates over the elements of the deque.
//
// The iterator panics if the deque has been modified since iteration started.
func (d *Deque[T]) Iterate() iterator.Iterator[T] {
	return &dequeIterator[T]{
		d:    d,
		i:    d.front,
		done: false,
		gen:  d.gen,
	}
}
