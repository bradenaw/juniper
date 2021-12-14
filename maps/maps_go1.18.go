//go:build go1.18

package maps

// Keys returns the keys of m as a slice.
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K,0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Values returns the values of m as a slice.
func Values[K comparable, V any](m map[K]V) []V {
	out := make([]V,0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}
