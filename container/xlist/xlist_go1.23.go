//go:build go1.23

package xlist

import (
	"iter"
)

func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for curr := l.front; curr != nil; curr = curr.next {
			if !yield(curr.Value) {
				break
			}
		}
	}
}

func (l *List[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		for curr := l.back; curr != nil; curr = curr.prev {
			if !yield(curr.Value) {
				break
			}
		}
	}
}

func (l *List[T]) Nodes() iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		curr := l.front
		for curr != nil {
			// Record next so that l.Remove() during iteration works.
			next := curr.next
			if !yield(curr) {
				break
			}
			curr = next
		}
	}
}
