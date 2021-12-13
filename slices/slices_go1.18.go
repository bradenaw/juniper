//go:build go1.18

package slices

// Grow grows x's capacity, if necessary, to fit n more elements and returns the modified slice.
// This does not change the length of x. After Grow(x, n), the following n append()s to x will not
// need to reallocate.
func Grow[T any](x []T, n int) []T {
	if cap(x)-len(x) < n {
		x2 := make([]T, len(x)+n)
		copy(x2, x)
		return x2[:len(x)]
	}
	return x
}

// Filter filters the contents of x to only those for which keep() returns true. This is done
// in-place and so modifies the contents of x. The modified slice is returned.
func Filter[T any](x []T, keep func(t T) bool) []T {
	filtered := x[:0]
	for i := range x {
		if keep(x[i]) {
			filtered = append(filtered, x[i])
		}
	}
	// Zero out the rest in case they contain pointers, so that filtered doesn't retain references.
	Clear(x[len(filtered):])
	return filtered
}

// Reverse reverses the elements of x in place.
func Reverse[T any](x []T) {
	for i := 0; i < len(x)/2; i++ {
		x[i], x[len(x)-i-1] = x[len(x)-i-1], x[i]
	}
}

// Insert inserts the given values starting at index idx, shifting elements after idx to the right
// and growing the slice to make room. Insert will expand the length of the slice up to its capacity
// if it can, if this isn't desired then x should be resliced to have capacity equal to its length:
//
//   x[:len(x):len(x)]
//
// The cost is linear in the number of elements added and the number of elements after idx that must
// be shifted.
func Insert[T any](x []T, idx int, values ...T) []T {
	x = Grow(x, len(values))
	x = x[: len(x)+len(values) : len(x)+len(values)]
	copy(x[idx+len(values):], x[idx:])
	copy(x[idx:], values)
	return x
}

// Remove removes n elements from x starting at index idx and returns the modified slice. This
// requires shifting the elements after the removed elements over, and so its cost is linear in the
// number of elements shifted.
func Remove[T any](x []T, idx int, n int) []T {
	copy(x[idx:], x[idx+n:])
	Clear(x[len(x)-n:])
	return x[:len(x)-n]
}

// Clear fills x with the zero value of T.
func Clear[T any](x []T) {
	var zero T
	for i := range x {
		x[i] = zero
	}
}

// Clone creates a new slice and copies the elements of x into it.
func Clone[T any](x []T) []T {
	return append([]T{}, x...)
}

// Compact removes adjacent duplicates from x in-place and returns the modified slice.
func Compact[T comparable](x []T) []T {
	compacted := x[:0]
	for i := range x {
		if i == 0 || x[i-1] != x[i] {
			compacted = append(compacted, x[i])
		}
	}
	Clear(x[len(compacted):])
	return compacted
}
