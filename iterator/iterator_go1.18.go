//go:build go1.18

// package iterator allows iterating over sequences of values, for example the contents of a
// container.
package iterator

// Iterator is used to iterate over a sequence of values.
//
// Iterators are lazy, meaning they do no work until a call to Next().
type Iterator[T any] interface {
	// Next advances the iterator and returns the next item. Once the iterator is finished, the
	// first return is meaningless and the second return is false. The final value of the iterator
	// will have true in the second return.
	Next() (T, bool)
}

type iterator[T any] struct {
	next func() (T, bool)
}

func (iter *iterator[T]) Next() (T, bool) {
	return iter.next()
}

// FromNext returns an iterator using the given next.
func FromNext[T any](next func() (T, bool)) Iterator[T] {
	return &iterator[T]{next: next}
}

type sliceIterator[T any] struct {
	a []T
}

func (iter *sliceIterator[T]) Next() (T, bool) {
	if len(iter.a) == 0 {
		var zero T
		return zero, false
	}
	item := iter.a[0]
	iter.a = iter.a[1:]
	return item, true
}

// Slice returns an iterator over the elements of s.
func Slice[T any](s []T) Iterator[T] {
	return &sliceIterator[T]{
		a: s,
	}
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

// Map transforms the values of iter using the conversion f.
type mapIterator[T any, U any] struct {
	inner Iterator[T]
	f     func(T) U
}

func (iter *mapIterator[T, U]) Next() (U, bool) {
	var zero U
	item, ok := iter.inner.Next()
	if !ok {
		return zero, false
	}
	return iter.f(item), true
}

// Map transforms the results of iter using the conversion f.
func Map[T any, U any](iter Iterator[T], f func(t T) U) Iterator[U] {
	return &mapIterator[T, U]{
		inner: iter,
		f:     f,
	}
}

type chunkIterator[T any] struct {
	chunkSize int
	inner     Iterator[T]
	chunk     []T
}

func (iter *chunkIterator[T]) Next() ([]T, bool) {
	chunk := make([]T, 0, iter.chunkSize)
	for {
		item, ok := iter.inner.Next()
		if !ok {
			break
		}
		chunk = append(chunk, item)
		if len(chunk) == iter.chunkSize {
			return chunk, true
		}
	}
	if len(chunk) > 0 {
		return chunk, true
	}
	return nil, false
}

// Chunk returns an iterator over non-overlapping chunks of size chunkSize. The last chunk will be
// smaller than chunkSize if the iterator does not contain an even multiple.
func Chunk[T any](iter Iterator[T], chunkSize int) Iterator[[]T] {
	return &chunkIterator[T]{
		inner:     iter,
		chunkSize: chunkSize,
	}
}

// Chain returns an Iterator that yields all elements of iters[0], then all elements of iters[1],
type chainIterator[T any] struct {
	iters []Iterator[T]
}

func (iter *chainIterator[T]) Next() (T, bool) {
	for len(iter.iters) > 0 {
		item, ok := iter.iters[0].Next()
		if ok {
			return item, true
		}
		iter.iters = iter.iters[1:]
	}
	var zero T
	return zero, false
}

// Chain returns an Iterator that returns all elements of iters[0], then all elements of iters[1],
// and so on.
func Chain[T any](iters ...Iterator[T]) Iterator[T] {
	return &chainIterator[T]{
		iters: iters,
	}
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

type filterIterator[T any] struct {
	inner Iterator[T]
	keep  func(T) bool
}

func (iter *filterIterator[T]) Next() (T, bool) {
	for {
		item, ok := iter.inner.Next()
		if !ok {
			break
		}
		if iter.keep(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// Filter returns an iterator that yields only the items from iter for which keep returns true.
func Filter[T any](iter Iterator[T], keep func(T) bool) Iterator[T] {
	return &filterIterator[T]{inner: iter, keep: keep}
}

type firstIterator[T any] struct {
	inner Iterator[T]
	x     int
}

func (iter *firstIterator[T]) Next() (T, bool) {
	if iter.x <= 0 {
		var zero T
		return zero, false
	}
	iter.x--
	return iter.inner.Next()
}

// First returns an iterator that yields the first n items from iter.
func First[T any](iter Iterator[T], n int) Iterator[T] {
	return &firstIterator[T]{inner: iter, x: n}
}

type whileIterator[T any] struct {
	inner Iterator[T]
	f     func(T) bool
	done  bool
}

func (iter *whileIterator[T]) Next() (T, bool) {
	var zero T
	if iter.done {
		return zero, false
	}
	item, ok := iter.Next()
	if !ok {
		return zero, false
	}
	if !iter.f(item) {
		iter.done = true
		return zero, false
	}
	return item, true
}

// While returns an iterator that terminates before the first item from iter for which f returns
// false.
func While[T any](iter Iterator[T], f func(T) bool) Iterator[T] {
	return &whileIterator[T]{
		inner: iter,
		f:     f,
		done:  false,
	}
}

type chanIterator[T any] struct {
	c <-chan T
}

func (iter *chanIterator[T]) Next() (T, bool) {
	item, ok := <-iter.c
	return item, ok
}

// Chan returns an Iterator that yields the values received on c.
func Chan[T any](c <-chan T) Iterator[T] {
	return &chanIterator[T]{c: c}
}
