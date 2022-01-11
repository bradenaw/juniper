//go:build go1.18

// package sets contains set operations like union, intersection, and difference.
package sets

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

type Set[T any] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Len() int
	Iterate() iterator.Iterator[T]
}

// Union adds to out out all items from sets and returns out.
func Union[T any](out Set[T], sets ...Set[T]) Set[T] {
	for _, set := range sets {
		iter := set.Iterate()
		for {
			item, ok := iter.Next()
			if !ok {
				break
			}
			out.Add(item)
		}
	}
	return out
}

// Intersection adds to out all items that appear in all sets and returns out.
func Intersection[T comparable](out Set[T], sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return out
	}

	// Ideally we check from most-selective to least-selective so we can do the fewest iterations
	// of each of the below loops. Use set size as an approximation.
	xsort.Slice(sets, func(a, b Set[T]) bool { return a.Len() < b.Len() })

	iter := sets[0].Iterate()
Outer:
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}

		for j := 1; j < len(sets); j++ {
			if !sets[j].Contains(item) {
				continue Outer
			}
		}
		out.Add(item)
	}
	return out
}

// Difference adds to out all items that appear in a but not in b and returns out.
func Difference[T comparable](out, a, b Set[T]) Set[T] {
	iter := a.Iterate()
	for {
		item, ok := iter.Next()
		if !ok {
			break
		}
		if !b.Contains(item) {
			out.Add(item)
		}
	}
	return out
}