// Package xslices contains utilities for working with slices of arbitrary types.
package xslices

// All returns true if f(s[i]) returns true for all i. Trivially, returns true if s is empty.
func All[T any](s []T, f func(T) bool) bool {
	for i := range s {
		if !f(s[i]) {
			return false
		}
	}
	return true
}

// Chunk returns non-overlapping chunks of s. The last chunk will be smaller than chunkSize if
// len(s) is not a multiple of chunkSize.
//
// Returns an empty slice if len(s)==0. Panics if chunkSize <= 0.
func Chunk[T any](s []T, chunkSize int) [][]T {
	out := make([][]T, (len(s)+chunkSize-1)/chunkSize)
	for i := range out {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(s) {
			end = len(s)
		}
		out[i] = s[start:end]
	}
	return out
}

// Clear fills s with the zero value of T.
//
// Deprecated: clear is a builtin as of Go 1.21.
func Clear[T any](s []T) {
	var zero T
	Fill(s, zero)
}

// Count returns the number of times x appears in s.
func Count[T comparable](s []T, x T) int {
	return CountFunc(s, func(s T) bool { return x == s })
}

// Count returns the number of items in s for which f returns true.
func CountFunc[T any](s []T, f func(T) bool) int {
	n := 0
	for _, s := range s {
		if f(s) {
			n++
		}
	}
	return n
}

// Fill fills s with copies of x.
func Fill[T any](s []T, x T) {
	for i := range s {
		s[i] = x
	}
}

// Group returns a map from u to all items of s for which f(s[i]) returned u.
func Group[T any, U comparable](s []T, f func(T) U) map[U][]T {
	m := make(map[U][]T)
	for i := range s {
		g := f(s[i])
		m[g] = append(m[g], s[i])
	}
	return m
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

// LastIndex returns the last index of x in s, or -1 if x is not in s.
func LastIndex[T comparable](s []T, x T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == x {
			return i
		}
	}
	return -1
}

// LastIndexFunc returns the last index in s for which f(s[i]) returns true, or -1 if there are no
// such items.
func LastIndexFunc[T any](s []T, f func(T) bool) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

// Map creates a new slice by applying f to each element of s.
func Map[T any, U any](s []T, f func(T) U) []U {
	out := make([]U, len(s))
	for i := range s {
		out[i] = f(s[i])
	}
	return out
}

// Partition moves elements of s such that all elements for which f returns false are at the
// beginning and all elements for which f returns true are at the end. It makes no other guarantees
// about the final order of elements. Returns the index of the first element for which f returned
// true, or len(s) if there wasn't one.
func Partition[T any](s []T, f func(t T) bool) int {
	i := 0
	j := len(s) - 1
	for {
		for i < j {
			if !f(s[i]) {
				i++
			} else {
				break
			}
		}
		for j > i {
			if f(s[j]) {
				j--
			} else {
				break
			}
		}
		if i >= j {
			break
		}
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
	if i < len(s) && !f(s[i]) {
		i++
	}
	return i
}

// Reduce reduces s to a single value using the reduction function f.
func Reduce[T any, U any](s []T, initial U, f func(U, T) U) U {
	out := initial
	for i := range s {
		out = f(out, s[i])
	}
	return out
}

// RemoveUnordered removes n elements from s starting at index idx and returns the modified slice.
// This is done by moving up to n elements from the end of the slice into the gap left by removal,
// which is linear in n (rather than len(s)-idx as Remove() is), but does not preserve order of the
// remaining elements.
func RemoveUnordered[T any](s []T, idx int, n int) []T {
	keepStart := len(s) - n
	removeEnd := idx + n
	if removeEnd > keepStart {
		keepStart = removeEnd
	}
	copy(s[idx:], s[keepStart:])
	Clear(s[len(s)-n:])
	return s[:len(s)-n]
}

// Repeat returns a slice with length n where every item is s.
func Repeat[T any](s T, n int) []T {
	out := make([]T, n)
	for i := range out {
		out[i] = s
	}
	return out
}

// Reverse reverses the elements of s in place.
func Reverse[T any](s []T) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

// Runs returns a slice of slices. The inner slices are contiguous runs of elements from s such that
// same(a, b) returns true for any a and b in the run.
//
// same(a, a) must return true. If same(a, b) and same(b, c) both return true, then same(a, c) must
// also.
//
// The returned slices use the same underlying array as s.
func Runs[T any](s []T, same func(a, b T) bool) [][]T {
	var runs [][]T
	start := 0
	end := 0
	for i := 1; i < len(s); i++ {
		if same(s[i-1], s[i]) {
			end = i + 1
		} else {
			runs = append(runs, s[start:end])
			start = i
			end = i + 1
		}
	}
	if end > 0 {
		runs = append(runs, s[start:])
	}
	return runs
}

// Shrink shrinks s's capacity by reallocating, if necessary, so that cap(s) <= len(s) + n.
func Shrink[T any](s []T, n int) []T {
	if cap(s) > len(s)+n {
		x2 := make([]T, len(s)+n)
		copy(x2, s)
		return x2[:len(s)]
	}
	return s
}

// Unique returns a slice that contains only the first instance of each unique item in s, preserving
// order.
//
// Compact is more efficient if duplicates are already adjacent in s, for example if s is in sorted
// order.
func Unique[T comparable](s []T) []T {
	return uniqueInto([]T{}, s)
}

// UniqueInPlace returns a slice that contains only the first instance of each unique item in s,
// preserving order. This is done in-place and so modifies the contents of s. The modified slice is
// returned.
//
// Compact is more efficient if duplicates are already adjacent in s, for example if s is in sorted
// order.
func UniqueInPlace[T comparable](s []T) []T {
	filtered := uniqueInto(s[:0], s)
	Clear(s[len(filtered):])
	return filtered
}

func uniqueInto[T comparable](into []T, s []T) []T {
	m := make(map[T]struct{}, len(s))
	for i := range s {
		_, ok := m[s[i]]
		if !ok {
			into = append(into, s[i])
			m[s[i]] = struct{}{}
		}
	}
	return into
}
