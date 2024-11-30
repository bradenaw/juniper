//go:build go1.23

package tree

import (
	"iter"
)

func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		it := s.Iterate()
		for {
			item, ok := it.Next()
			if !ok {
				return
			}
			if !yield(item) {
				return
			}
		}
	}
}

func (s Set[T]) Backward() iter.Seq[T] {
	return func(yield func(T) bool) {
		it := s.RangeReverse(Unbounded[T](), Unbounded[T]())
		for {
			item, ok := it.Next()
			if !ok {
				return
			}
			if !yield(item) {
				return
			}
		}
	}
}
