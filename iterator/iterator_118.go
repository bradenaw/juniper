//+build go1.18

package iterator

type Iterator[T any] struct {
	inner func() (T, bool)
	item T
}

func (iter *Iterator[T]) Next() bool {
	item, ok := iter.inner()
	iter.item = item
	return ok
}

func (iter *Iterator[T]) Item() T {
	return iter.item
}

func (iter *Iterator[T]) Collect() []T {
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
