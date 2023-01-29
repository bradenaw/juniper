package maps

import (
	"github.com/bradenaw/juniper/xslices"
	"github.com/bradenaw/juniper/xsort"
)

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func Clone[M ~map[K]V, K comparable, V any](m M) map[K]V {
	result := make(map[K]V, len(m))
	Copy(result, m)
	return result
}

func Copy[M ~map[K]V, K comparable, V any](out, m M) {
	for k, v := range m {
		out[k] = v
	}
}

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

func FromKeys[K comparable](ks []K) map[K]struct{} {
	result := make(map[K]struct{}, len(ks))
	for _, k := range ks {
		result[k] = struct{}{}
	}
	return result
}

func Union[T comparable](sets ...map[T]struct{}) map[T]struct{} {
	size := 0
	first := true
	for _, set := range sets {
		if first || len(set) < size {
			size = len(set)
		}
		first = false
	}
	out := make(map[T]struct{}, size)
	for _, set := range sets {
		for k := range set {
			out[k] = struct{}{}
		}
	}
	return out
}

func Intersection[T comparable](sets ...map[T]struct{}) map[T]struct{} {
	out := make(map[T]struct{})
	if len(sets) == 0 {
		return out
	}

	xsort.Slice(xslices.Clone(sets), func(a, b map[T]struct{}) bool { return len(a) < len(b) })

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

func Intersects[T comparable](sets ...map[T]struct{}) bool {
	if len(sets) == 0 {
		return false
	}

	// Ideally we check from most-selective to least-selective so we can do the fewest iterations
	// of each of the below loops. Use set size as an approximation.
	xsort.Slice(xslices.Clone(sets), func(a, b map[T]struct{}) bool { return len(a) < len(b) })

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

func Difference[T comparable](a, b map[T]struct{}) map[T]struct{} {
	size := len(a) - len(b)
	if size < 0 {
		size = 0
	}
	result := make(map[T]struct{}, size)
	for k := range a {
		if _, ok := b[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}
