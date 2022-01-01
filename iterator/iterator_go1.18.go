//go:build go1.18

// package iterator allows iterating over sequences of values, for example the contents of a
// container.
package iterator

// Iterator is used to iterate over a sequence of values.
//
// Iterators are lazy, meaning they do no work until a call to Next().
type Iterator[T any] interface {
	// Next advances the iterator to the next item in the sequence. Returns false if the iterator
	// is now past the end of the sequence.
	Next() (T, bool)
}

type iterator[T any] struct {
	next func() (T, bool)
}

// FromNext returns an iterator using the given next.
func FromNext[T any](next func() (T, bool)) Iterator[T] {
	return &iterator[T]{next: next}
}

func (iter *iterator[T]) Next() (T, bool) {
	return iter.next()
}

// Slice returns an iterator over the elements of s.
func Slice[T any](s []T) Iterator[T] {
	i := 0
	return FromNext(func() (T, bool) {
		if i >= len(s) {
			var zero T
			return zero, false
		}
		item := s[i]
		i++
		return item, true
	})
}

// Collect advances iter to the end and returns all of the items seen as a slice.
func Collect[T any](iter Iterator[T]) []T {
	var out []T
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		out = append(out, item)
	}
	return out
}

// Map transforms the results of iter using the conversion f.
func Map[T any, U any](iter Iterator[T], f func(t T) U) Iterator[U] {
	return FromNext(func() (U, bool) {
		var zero U
		item, ok := iter.Next()
		if !ok {
			return zero, false
		}
		return f(item), true
	})
}

// Chunk returns an iterator over non-overlapping chunks of size chunkSize. The last chunk will be
// smaller than chunkSize if the iterator does not contain an even multiple.
func Chunk[T any](iter Iterator[T], chunkSize int) Iterator[[]T] {
	return FromNext(func() ([]T, bool) {
		chunk := make([]T, 0, chunkSize)
		for {
			item, ok := iter.Next()
			if !ok {
				break
			}
			chunk = append(chunk, item)
			if len(chunk) == chunkSize {
				return chunk, true
			}
		}
		if len(chunk) > 0 {
			return chunk, true
		}
		return nil, false
	})
}

// Chain returns an Iterator that returns all elements of iters[0], then all elements of iters[1],
// and so on.
func Chain[T any](iters ...Iterator[T]) Iterator[T] {
	i := 0
	return FromNext(func() (T, bool) {
		var zero T
		for {
			if i >= len(iters) {
				return zero, false
			}
			item, ok := iters[i].Next()
			if ok {
				return item, true
			}
			i++
		}
	})
}

// Equal returns true if the given iterators yield the same items in the same order. Consumes the
// iterators.
func Equal[T comparable](iters ...Iterator[T]) bool {
	if len(iters) == 0 {
		return true
	}
	for {
		item, ok := iters[0].Next()
		for i := 1; i < len(iters); i++ {
			iterIItem, iterIOk := iters[i].Next()
			if ok != iterIOk {
				return false
			}
			if item != iterIItem {
				return false
			}
		}
		if !ok {
			return true
		}
	}
}

// Filter returns an iterator that yields only the items from iter for which keep returns true.
func Filter[T any](iter Iterator[T], keep func(T) bool) Iterator[T] {
	return FromNext(func() (T, bool) {
		for {
			item, ok := iter.Next()
			if !ok {
				break
			}
			if keep(item) {
				return item, true
			}
		}
		var zero T
		return zero, false
	})
}

// First returns an iterator that yields the first n items from iter.
func First[T any](iter Iterator[T], n int) Iterator[T] {
	i := 0
	return FromNext(func() (T, bool) {
		if i >= n {
			var zero T
			return zero, false
		}
		i++
		return iter.Next()
	})
}

// While returns an iterator that terminates at the first item from iter for which f returns false.
func While[T any](iter Iterator[T], f func(T) bool) Iterator[T] {
	done := false
	return FromNext(func() (T, bool) {
		var zero T
		if done {
			return zero, false
		}
		item, ok := iter.Next()
		if !ok {
			return zero, false
		}
		if !f(item) {
			done = true
			return zero, false
		}
		return item, true
	})
}
