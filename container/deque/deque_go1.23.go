//go:build go1.23

package deque

import (
	"iter"
)

func (d *Deque[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < d.Len(); i++ {
			idx := (d.front + i) % len(d.a)
			if !yield(d.a[idx]) {
				return
			}
		}
	}
}

func (d *Deque[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := d.Len() - 1; i >= 0; i-- {
			idx := (d.front + i) % len(d.a)
			if !yield(d.a[idx]) {
				return
			}
		}
	}
}
