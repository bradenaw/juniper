package xmaps

import (
	"github.com/bradenaw/juniper/xslices"
	"github.com/bradenaw/juniper/xsort"
)

func Reverse[M ~map[K]V, K comparable, V comparable](m M) map[V][]K {
	result := make(map[V][]K, len(m))
	for k, v := range m {
		result[v] = append(result[v], k)
	}
	return result
}

func ReverseSingle[M ~map[K]V, K comparable, V comparable](m M) (map[V]K, bool) {
	result := make(map[V]K, len(m))
	ok := true
	for k, v := range m {
		if _, ok := result[v]; ok {
			ok = false
		}
		result[v] = k
	}
	return result, ok
}

type Set[T comparable] map[T]struct{}

func SetFromSlice[T comparable](items []T) Set[T] {
	result := make(Set[T], len(items))
	for _, k := range items {
		result[k] = struct{}{}
	}
	return result
}

func Union[S ~map[T]struct{}, T comparable](sets ...S) S {
	size := 0
	first := true
	for _, set := range sets {
		if first || len(set) < size {
			size = len(set)
		}
		first = false
	}
	out := make(S, size)
	for _, set := range sets {
		for k := range set {
			out[k] = struct{}{}
		}
	}
	return out
}

func Intersection[S ~map[T]struct{}, T comparable](sets ...S) S {
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

func Difference[S ~map[T]struct{}, T comparable](a, b S) S {
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
