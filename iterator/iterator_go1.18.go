//go:build go1.18

package iterator

type Iterator[T any] struct {
	inner func() (T, bool)
	item  T
}

func New[T any](f func() (T, bool)) Iterator[T] {
	return Iterator[T]{inner: f}
}

func (iter *Iterator[T]) Next() bool {
	item, ok := iter.inner()
	iter.item = item
	return ok
}

func (iter *Iterator[T]) Item() T {
	return iter.item
}

func Slice[T any](s []T) Iterator[T] {
	i := 0
	return Iterator[T]{inner: func() (T, bool) {
		if i >= len(s) {
			var zero T
			return zero, false
		}
		item := s[i]
		i++
		return item, true
	}}
}

func Collect[T any](iter Iterator[T]) []T {
	var out []T
	for iter.Next() {
		out = append(out, iter.Item())
	}
	return out
}

func Map[T any, U any](iter Iterator[T], f func(t T) U) Iterator[U] {
	return Iterator[U]{
		inner: func() (U, bool) {
			var zero U
			if !iter.Next() {
				return zero, false
			}
			return f(iter.Item()), true
		},
	}
}

func Chunk[T any](iter Iterator[T], chunkSize int) Iterator[[]T] {
	chunk := make([]T, 0, chunkSize)
	return Iterator[[]T]{
		inner: func() ([]T, bool) {
			for iter.Next() {
				chunk = append(chunk, iter.Item())
				if len(chunk) == chunkSize {
					item := chunk
					chunk = chunk[:0]
					return item, true
				}
			}
			if len(chunk) > 0 {
				item := chunk
				chunk = chunk[:0]
				return item, true
			}
			return nil, false
		},
	}
}

func Chain[T any](iters ...Iterator[T]) Iterator[T] {
	i := 0
	return Iterator[T]{
		inner: func() (T, bool) {
			var zero T
			for {
				if i >= len(iters) {
					return zero, false
				}
				if iters[i].Next() {
					return iters[i].Item(), true
				}
				i++
			}
		},
	}
}
