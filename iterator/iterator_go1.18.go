//go:build go1.18

package iterator

// Iterator is used to iterate over a sequence of values.
//
// On creation, iterators are placed 'before' the sequence, and so Next() must be called to advance
// to the first item.
//
// Iterators are lazy, meaning they do no work until a call to Next().
type Iterator[T any] interface {
	// Next advances the iterator to the next item in the sequence. Returns false if the iterator
	// is now past the end of the sequence.
	Next() bool
	// Item returns the item the iterator is currently at.
	Item() T
}

type iterator[T any] struct {
	inner func() (T, bool)
	item  T
}

// New returns an iterator using a short-hand function next. The second return is false when next
// has advanced past the end of the sequence.
func New[T any](next func() (T, bool)) Iterator[T] {
	return &iterator[T]{inner: next}
}

func (iter *iterator[T]) Next() bool {
	item, ok := iter.inner()
	iter.item = item
	return ok
}

func (iter *iterator[T]) Item() T {
	return iter.item
}

// Slice returns an iterator over the elements of s.
func Slice[T any](s []T) Iterator[T] {
	i := 0
	return &iterator[T]{inner: func() (T, bool) {
		if i >= len(s) {
			var zero T
			return zero, false
		}
		item := s[i]
		i++
		return item, true
	}}
}

// Collect advances iter to the end and returns all of the items seen as a slice.
func Collect[T any](iter Iterator[T]) []T {
	var out []T
	for iter.Next() {
		out = append(out, iter.Item())
	}
	return out
}

// Map transforms the results of iter using the conversion f.
func Map[T any, U any](iter Iterator[T], f func(t T) U) Iterator[U] {
	return &iterator[U]{
		inner: func() (U, bool) {
			var zero U
			if !iter.Next() {
				return zero, false
			}
			return f(iter.Item()), true
		},
	}
}

// Chunk returns an iterator over non-overlapping chunks of size chunkSize. The last chunk will be
// smaller than chunkSize if the iterator does not contain an even multiple.
func Chunk[T any](iter Iterator[T], chunkSize int) Iterator[[]T] {
	chunk := make([]T, 0, chunkSize)
	return &iterator[[]T]{
		inner: func() ([]T, bool) {
			for iter.Next() {
				chunk = append(chunk, iter.Item())
				if len(chunk) == chunkSize {
					item := chunk
					chunk = make([]T, 0, chunkSize)
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

// Chain returns an Iterator that returns all elements of iters[0], then all elements of iters[1],
// and so on.
func Chain[T any](iters ...Iterator[T]) Iterator[T] {
	i := 0
	return &iterator[T]{
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
