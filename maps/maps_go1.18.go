//go:build go1.18

package maps

import (
	"reflect"

	"github.com/bradenaw/juniper/iterator"
)

// Keys returns the keys of m as a slice.
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Values returns the values of m as a slice.
func Values[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// Set implements sets.Set for map[T]struct{}.
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

type setIterator[T any] struct {
	// Unfortunately map range is special, so this appears to be the only way to do this. The only
	// alternatives are to (a) use a goroutine, which we'd leak if the caller didn't drain the
	// iterator, e.g. by using iterator.First(); or (b) to change the sets.Set interface to have
	// another method that takes a function closure the way that sync.Map does.
	//
	// (a) seems an obvious nonstarter.
	//
	// (b) would perform better but with a worse API. Stick the better API for now and hope that the
	// runtime grows a way to do this faster later on.
	inner *reflect.MapIter
}

func (iter *setIterator[T]) Next() (T, bool) {
	if iter.inner == nil || !iter.inner.Next() {
		var zero T
		iter.inner = nil
		return zero, false
	}
	return iter.inner.Key().Interface().(T), true
}

func (s Set[T]) Iterate() iterator.Iterator[T] {
	return &setIterator[T]{inner: reflect.ValueOf(s).MapRange()}
}
