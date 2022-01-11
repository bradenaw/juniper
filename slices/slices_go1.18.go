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
// The time cost is O(n+m) where n is len(values) and m is len(x[idx:]).
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

// Equal returns true if a and b contain the same items in the same order.
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Count returns the number of times item appears in a.
func Count[T comparable](a []T, item T) int {
	return CountFunc(a, func(x T) bool { return item == x })
}

// Count returns the number of items in a for which f returns true.
func CountFunc[T any](a []T, f func(T) bool) int {
	n := 0
	for _, x := range a {
		if f(x) {
			n++
		}
	}
	return n
}

// Index returns the first index of item in a, or -1 if item is not in a.
func Index[T comparable](a []T, item T) int {
	for i := range a {
		if a[i] == item {
			return i
		}
	}
	return -1
}

// LastIndex returns the last index of item in a, or -1 if item is not in a.
func LastIndex[T comparable](a []T, item T) int {
	for i := len(a) - 1; i >= 0; i-- {
		if a[i] == item {
			return i
		}
	}
	return -1
}

// Index returns the first index in a for which f(a[i]) returns true, or -1 if there are no such
// items.
func IndexFunc[T any](a []T, f func(T) bool) int {
	for i := range a {
		if f(a[i]) {
			return i
		}
	}
	return -1
}

// LastIndexFunc returns the last index in a for which f(a[i]) returns true, or -1 if there are no
// such items.
func LastIndexFunc[T any](a []T, f func(T) bool ) int {
	for i := len(a) - 1; i >= 0; i-- {
		if f(a[i]) {
			return i
		}
	}
	return -1
}

// Join joins together the contents of each in.
func Join[T any](in ...[]T) []T {
	n := 0
	for i := range in {
		n += len(in[i])
	}
	out := make([]T, 0, n)
	for i := range in {
		out = append(out, in[i]...)
	}
	return out
}

// Partition moves elements of x such that all elements for which f returns false are at the
// beginning and all elements for which f returns true are at the end. It makes no other guarantees
// about the final order of elements.
func Partition[T any](x []T, f func(t T) bool) {
	i := 0
	j := len(x) - 1
	for {
		for i < j {
			if !f(x[i]) {
				i++
			} else {
				break
			}
		}
		for j > i {
			if f(x[j]) {
				j--
			} else {
				break
			}
		}
		if i >= j {
			break
		}
		x[i], x[j] = x[j], x[i]
		i++
		j--
	}
}

// Unique removes duplicates from x in-place, preserving order, and returns the modified slice.
//
// Compact is more efficient if duplicates are already adjacent in x, for example if x is in sorted
// order.
func Unique[T comparable](x []T) []T {
	filtered := x[:0]
	m := make(map[T]struct{}, len(x))
	for i := range x {
		_, ok := m[x[i]]
		if !ok {
			filtered = append(filtered, x[i])
			m[x[i]] = struct{}{}
		}
	}
	Clear(x[len(filtered):])
	return filtered
}

// Chunk returns non-overlapping chunks of x. The last chunk will be smaller than chunkSize if
// len(x) is not a multiple of chunkSize.
//
// Returns an empty slice if len(x)==0. Panics if chunkSize <= 0.
func Chunk[T any](x []T, chunkSize int) [][]T {
	out := make([][]T, (len(x)+chunkSize-1)/chunkSize)
	for i := range out {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(x) {
			end = len(x)
		}
		out[i] = x[start:end]
	}
	return out
}

// Any returns true if f(x[i]) returns true for any i. Trivially, returns false if x is empty.
func Any[T any](x []T, f func(T) bool) bool {
	for i := range x {
		if f(x[i]) {
			return true
		}
	}
	return false
}

// All returns true if f(x[i]) returns true for all i. Trivially, returns true if x is empty.
func All[T any](x []T, f func(T) bool) bool {
	for i := range x {
		if !f(x[i]) {
			return false
		}
	}
	return true
}

// Map creates a new slice by applying f to each element of x.
func Map[T any, U any](x []T, f func(T) U) []U {
	out := make([]U, len(x))
	for i := range x {
		out[i] = f(x[i])
	}
	return out
}

// Runs returns a slice of slices. The inner slices are contiguous runs of elements from x such
// that same(a, b) returns true for any a and b in the run.
//
// same(a, a) must return true. If same(a, b) and same(b, c) both return true, then same(a, c) must
// also.
//
// The returned slices use the same underlying array as x.
func Runs[T any](x []T, same func(a, b T) bool) [][]T {
	var runs [][]T
	start := 0
	end := 0
	for i := 1; i < len(x); i++ {
		if same(x[i-1], x[i]) {
			end = i + 1
		} else {
			runs = append(runs, x[start:end])
			start = i
			end = i + 1
		}
	}
	if end > 0 {
		runs = append(runs, x[start:])
	}
	return runs
}
