// Package sets contains set operations like union, intersection, and difference.
package sets

import (
	"reflect"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

// Map implements sets.Set for map[T]struct{}.
type Map[T comparable] map[T]struct{}

func (s Map[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Map[T]) Remove(item T) {
	delete(s, item)
}

func (s Map[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Map[T]) Len() int {
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

func (s Map[T]) Iterate() iterator.Iterator[T] {
	return &setIterator[T]{inner: reflect.ValueOf(s).MapRange()}
}

func (s Map[T]) IterateInternal(f func(T) bool) {
	for t := range s {
		if !f(t) {
			break
		}
	}
}

var _ withIterateInternal[byte] = Map[byte]{}

// Set is a minimal interface to a set. It is implemented by sets.Map and container/tree.Set, among
// others.
type Set[T any] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Len() int
	Iterate() iterator.Iterator[T]
}

type withIterateInternal[T any] interface {
	IterateInternal(func(T) bool)
}

func iterateInternal[T any](s Set[T], f func(T) bool) {
	if s, ok := s.(withIterateInternal[T]); ok {
		s.IterateInternal(f)
		return
	}
	iter := s.Iterate()
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		if !f(item) {
			break
		}
	}
}

// Union adds to out out all items from sets and returns out.
func Union[T any](out Set[T], sets ...Set[T]) Set[T] {
	for _, set := range sets {
		iterateInternal(set, func(item T) bool {
			out.Add(item)
			return true
		})
	}
	return out
}

// Intersection adds to out all items that appear in all sets and returns out.
func Intersection[T any](out Set[T], sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return out
	}

	// Ideally we check from most-selective to least-selective so we can do the fewest iterations
	// of each of the below loops. Use set size as an approximation.
	xsort.Slice(sets, func(a, b Set[T]) bool { return a.Len() < b.Len() })

	iterateInternal(sets[0], func(item T) bool {
		for j := 1; j < len(sets); j++ {
			if !sets[j].Contains(item) {
				return true
			}
		}
		out.Add(item)
		return true
	})
	return out
}

// Intersects returns true if the given sets have any items in common.
func Intersects[T any](sets ...Set[T]) bool {
	if len(sets) == 0 {
		return false
	}

	// Ideally we check from most-selective to least-selective so we can do the fewest iterations
	// of each of the below loops. Use set size as an approximation.
	xsort.Slice(sets, func(a, b Set[T]) bool { return a.Len() < b.Len() })

	out := false
	iterateInternal(sets[0], func(item T) bool {
		for j := 1; j < len(sets); j++ {
			if !sets[j].Contains(item) {
				return true
			}
		}
		out = true
		return false
	})
	return out
}

// Difference adds to out all items that appear in a but not in b and returns out.
func Difference[T any](out, a, b Set[T]) Set[T] {
	iterateInternal(a, func(item T) bool {
		if !b.Contains(item) {
			out.Add(item)
		}
		return true
	})
	return out
}
