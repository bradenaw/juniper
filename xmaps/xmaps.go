// Package xmaps contains utilities for working with maps.
package xmaps

import (
	"fmt"

	"github.com/bradenaw/juniper/xslices"
	"github.com/bradenaw/juniper/xsort"
)

// Reverse returns a map from m's values to each of the keys that mapped to it in arbitrary order.
func Reverse[M ~map[K]V, K comparable, V comparable](m M) map[V][]K {
	result := make(map[V][]K, len(m))
	for k, v := range m {
		result[v] = append(result[v], k)
	}
	return result
}

// ReverseSingle returns a map of m's values to m's keys. If there are any duplicate values, the
// resulting map has an arbitrary choice of the associated keys and the second return is false.
func ReverseSingle[M ~map[K]V, K comparable, V comparable](m M) (map[V]K, bool) {
	result := make(map[V]K, len(m))
	allOk := true
	for k, v := range m {
		if _, ok := result[v]; ok {
			allOk = false
		}
		result[v] = k
	}
	return result, allOk
}

// ToIndex returns a map from keys[i] to i.
func ToIndex[K comparable](keys []K) map[K]int {
	m := make(map[K]int, len(keys))
	for i := range keys {
		m[keys[i]] = i
	}
	return m
}

// FromKeysAndValues returns a map from keys[i] to values[i]. If there are any duplicate keys, the
// resulting map has an arbitrary choice of the associated values and the second return is false. It
// panics if len(keys)!=len(values).
func FromKeysAndValues[K comparable, V any](keys []K, values []V) (map[K]V, bool) {
	if len(keys) != len(values) {
		panic(fmt.Sprintf("len(keys)=%d, len(values)=%d", len(keys), len(values)))
	}
	m := make(map[K]V, len(keys))
	allOk := true
	for i := range keys {
		if _, ok := m[keys[i]]; ok {
			allOk = false
		}
		m[keys[i]] = values[i]
	}
	return m, allOk
}

// Set[T] is shorthand for map[T]struct{} with convenience methods.
type Set[T comparable] map[T]struct{}

// Add adds item to the set.
func (s Set[T]) Add(item T) { s[item] = struct{}{} }

// Remove removes item from the set.
func (s Set[T]) Remove(item T) { delete(s, item) }

// Contains returns true if item is in the set.
func (s Set[T]) Contains(item T) bool { _, ok := s[item]; return ok }

// SetFromSlice returns a Set whose elements are items.
func SetFromSlice[T comparable](items []T) Set[T] {
	result := make(Set[T], len(items))
	for _, k := range items {
		result[k] = struct{}{}
	}
	return result
}

// Union returns a set containing all elements of all input sets.
func Union[S ~map[T]struct{}, T comparable](sets ...S) S {
	// Size estimate: the smallest possible result is the largest input set, if it's a superset of
	// all of the others.
	size := 0
	for _, set := range sets {
		if len(set) > size {
			size = len(set)
		}
	}
	out := make(S, size)

	for _, set := range sets {
		for k := range set {
			out[k] = struct{}{}
		}
	}
	return out
}

// Intersection returns a set of the items that all input sets have in common.
func Intersection[S ~map[T]struct{}, T comparable](sets ...S) S {
	// The smallest intersection is 0, so don't guess about capacity.
	out := make(S)
	if len(sets) == 0 {
		return out
	}

	xsort.Slice(xslices.Clone(sets), func(a, b S) bool { return len(a) < len(b) })

	for k := range sets[0] {
		include := true
		for j := 1; j < len(sets); j++ {
			if _, ok := sets[j][k]; !ok {
				include = false
				break
			}
		}
		if include {
			out[k] = struct{}{}
		}
	}
	return out
}

// Intersects returns true if the input sets have any element in common.
func Intersects[S ~map[T]struct{}, T comparable](sets ...S) bool {
	if len(sets) == 0 {
		return false
	}

	// Ideally we check from most-selective to least-selective so we can do the fewest iterations
	// of each of the below loops. Use set size as an approximation.
	xsort.Slice(xslices.Clone(sets), func(a, b S) bool { return len(a) < len(b) })

	for k := range sets[0] {
		include := true
		for j := 1; j < len(sets); j++ {
			if _, ok := sets[j][k]; !ok {
				include = false
				break
			}
		}
		if include {
			return true
		}
	}
	return false
}

// Difference returns all items of a that do not appear in b.
func Difference[S ~map[T]struct{}, T comparable](a, b S) S {
	// Size estimate: the smallest possible result is if all items of b are in a.
	size := len(a) - len(b)
	if size < 0 {
		size = 0
	}
	result := make(S, size)

	for k := range a {
		if _, ok := b[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}
