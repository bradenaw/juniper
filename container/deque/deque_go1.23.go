//go:build go1.23

package deque

import (
	"iter"
)

func (d *Deque[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range d.Len() {
			if !yield(i, d.Item(i)) {
				return
			}
		}
	}
}

func (d *Deque[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := d.Len() - 1; i >= 0; i-- {
			if !yield(i, d.Item(i)) {
				return
			}
		}
	}
}

func (d *Deque[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range d.Len() {
			if !yield(d.Item(i)) {
				return
			}
		}
	}
}
