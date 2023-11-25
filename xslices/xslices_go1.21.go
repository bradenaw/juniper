//go:build go1.21

package xslices

import (
	"slices"
)

// Any returns true if f(s[i]) returns true for any i. Trivially, returns false if s is empty.
//
// Deprecated: slices.ContainsFunc is in the standard library as of Go 1.21.
func Any[T any](s []T, f func(T) bool) bool {
	return slices.ContainsFunc(s, f)
}

// Clone creates a new slice and copies the elements of s into it.
//
// Deprecated: slices.Clone is in the standard library as of Go 1.21.
func Clone[T any](s []T) []T {
	return slices.Clone(s)
}

// Compact returns a slice containing only the first item from each contiguous run of the same item.
//
// For example, this can be used to remove duplicates more cheaply than Unique when the slice is
// already in sorted order.
//
// Deprecated: slices.Compact(slices.Clone(s)) is in the standard library as of Go 1.21.
func Compact[T comparable](s []T) []T {
	return slices.Compact(slices.Clone(s))
}

// CompactInPlace returns a slice containing only the first item from each contiguous run of the
// same item. This is done in-place and so modifies the contents of s. The modified slice is
// returned.
//
// For example, this can be used to remove duplicates more cheaply than Unique when the slice is
// already in sorted order.
//
// Deprecated: slices.Compact is in the standard library as of Go 1.21.
func CompactInPlace[T comparable](s []T) []T {
	return slices.Compact(s)
}

// CompactFunc returns a slice containing only the first item from each contiguous run of items for
// which eq returns true.
//
// Deprecated: slices.CompactFunc(slices.Clone(s)) is in the standard library as of Go 1.21.
func CompactFunc[T any](s []T, eq func(T, T) bool) []T {
	return slices.CompactFunc(slices.Clone(s), eq)
}

// CompactInPlaceFunc returns a slice containing only the first item from each contiguous run of
// items for which eq returns true. This is done in-place and so modifies the contents of s. The
// modified slice is returned.
//
// Deprecated: slices.CompactFunc is in the standard library as of Go 1.21.
func CompactInPlaceFunc[T any](s []T, eq func(T, T) bool) []T {
	return slices.CompactFunc(s, eq)
}

// Equal returns true if a and b contain the same items in the same order.
//
// Deprecated: slices.Equal is in the standard library as of Go 1.21.
func Equal[T comparable](a, b []T) bool {
	return slices.Equal(a, b)
}

// EqualFunc returns true if a and b contain the same items in the same order according to eq.
//
// Deprecated: slices.EqualFunc is in the standard library as of Go 1.21.
func EqualFunc[T any](a, b []T, eq func(T, T) bool) bool {
	return slices.EqualFunc(a, b, eq)
}

// Filter returns a slice containing only the elements of s for which keep() returns true in the
// same order that they appeared in s.
//
// Deprecated: slices.DeleteFunc(slices.Clone(s), f) is in the standard library as of Go 1.21,
// though the polarity of the passed function is opposite: return true to remove, rather than to
// retain.
func Filter[T any](s []T, keep func(t T) bool) []T {
	return slices.DeleteFunc(slices.Clone(s), func(t T) bool { return !keep(t) })
}

// FilterInPlace returns a slice containing only the elements of s for which keep() returns true in
// the same order that they appeared in s. This is done in-place and so modifies the contents of s.
// The modified slice is returned.
//
// Deprecated: slices.DeleteFunc is in the standard library as of Go 1.21, though the polarity of
// the passed function is opposite: return true to remove, rather than to retain.
func FilterInPlace[T any](s []T, keep func(t T) bool) []T {
	return slices.DeleteFunc(s, func(t T) bool { return !keep(t) })
}

// Grow grows s's capacity by reallocating, if necessary, to fit n more elements and returns the
// modified slice. This does not change the length of s. After Grow(s, n), the following n
// append()s to s will not need to reallocate.
//
// Deprecated: slices.Grow is in the standard library as of Go 1.21.
func Grow[T any](s []T, n int) []T {
	return slices.Grow(s, n)
}

// Index returns the first index of x in s, or -1 if x is not in s.
//
// Deprecated: slices.Index is in the standard library as of Go 1.21.
func Index[T comparable](s []T, x T) int {
	return slices.Index(s, x)
}

// Index returns the first index in s for which f(s[i]) returns true, or -1 if there are no such
// items.
//
// Deprecated: slices.IndexFunc is in the standard library as of Go 1.21.
func IndexFunc[T any](s []T, f func(T) bool) int {
	return slices.IndexFunc(s, f)
}

// Insert inserts the given values starting at index idx, shifting elements after idx to the right
// and growing the slice to make room. Insert will expand the length of the slice up to its capacity
// if it can, if this isn't desired then s should be resliced to have capacity equal to its length:
//
//	s[:len(s):len(s)]
//
// The time cost is O(n+m) where n is len(values) and m is len(s[idx:]).
//
// Deprecated: slices.Insert is in the standard library as of Go 1.21.
func Insert[T any](s []T, idx int, values ...T) []T {
	return slices.Insert(s, idx, values...)
}

// Remove removes n elements from s starting at index idx and returns the modified slice. This
// requires shifting the elements after the removed elements over, and so its cost is linear in the
// number of elements shifted.
//
// Deprecated: slices.Delete is in the standard library as of Go 1.21, though slices.Delete takes
// two indexes rather than an index and a length.
func Remove[T any](s []T, idx int, n int) []T {
	return slices.Delete(s, idx, idx+n)
}
