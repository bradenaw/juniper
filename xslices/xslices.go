// Package xslices contains utilities for working with slices of arbitrary types.
package xslices

import (
	"cmp"
	"slices"

	"github.com/bradenaw/juniper/internal/heap"
)

// All returns true if f(s[i]) returns true for all i. Trivially, returns true if s is empty.
func All[T any](s []T, f func(T) bool) bool {
	for i := range s {
		if !f(s[i]) {
			return false
		}
	}
	return true
}

// Any returns true if f(s[i]) returns true for any i. Trivially, returns false if s is empty.
//
// Deprecated: slices.ContainsFunc is in the standard library as of Go 1.21.
func Any[T any](s []T, f func(T) bool) bool {
	return slices.ContainsFunc(s, f)
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

// Fill fills s with copies of x.
func Fill[T any](s []T, x T) {
	for i := range s {
		s[i] = x
	}
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

// Group returns a map from u to all items of s for which f(s[i]) returned u.
func Group[T any, U comparable](s []T, f func(T) U) map[U][]T {
	m := make(map[U][]T)
	for i := range s {
		g := f(s[i])
		m[g] = append(m[g], s[i])
	}
	return m
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

// MaxK returns the k maximum items from s in descending order. If s has fewer than k items, MaxK
// returns all of them.
func MaxK[S ~[]E, E cmp.Ordered](s S, k int) []E {
	return MaxKFunc(s, k, cmp.Compare[E])
}

// MaxKFunc returns the k maximum items according to cmp from s in descending order. If s has fewer
// than k items, MaxKFunc returns all of them.
func MaxKFunc[S ~[]E, E any](s S, k int, cmp func(E, E) int) []E {
	h := heap.New[E](cmp, func(a E, i int) {}, nil)
	h.Grow(k)

	for i := range s {
		h.Push(s[i])
		if h.Len() > k {
			h.Pop()
		}
	}
	out := make([]E, h.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = h.Pop()
	}
	return out
}

// MinK returns the k minimum items from s in ascending order. If s has fewer than k items, MinK
// returns all of them.
func MinK[S ~[]E, E cmp.Ordered](s S, k int) []E {
	return MinKFunc(s, k, cmp.Compare[E])
}

// MinKFunc returns the k minimum items according to cmp from s in ascending order. If s has fewer
// than k items, MinKFunc returns all of them.
func MinKFunc[S ~[]E, E any](s S, k int, cmp func(E, E) int) []E {
	return MaxKFunc(s, k, func(a, b E) int {
		// Copy of xsort.ReverseCompare to avoid dependency cycle.
		c := cmp(a, b)
		if -c == c {
			return 1
		}
		return -c
	})
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

// Remove removes n elements from s starting at index idx and returns the modified slice. This
// requires shifting the elements after the removed elements over, and so its cost is linear in the
// number of elements shifted.
//
// Deprecated: slices.Delete is in the standard library as of Go 1.21, though slices.Delete takes
// two indexes rather than an index and a length.
func Remove[T any](s []T, idx int, n int) []T {
	return slices.Delete(s, idx, idx+n)
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

func SortedMerge[T cmp.Ordered](in ...[]T) []T {
	return SortedMergeFunc(in, cmp.Compare[T])
}

func SortedMerge2[T cmp.Ordered](a []T, b []T) []T {
	return SortedMerge2Func(a, b, cmp.Compare[T])
}

type valueAndSource[T any] struct {
	value  T
	source int
}

func SortedMerge2Func[T any](a []T, b []T, cmp func(T, T) int) []T {
	out := make([]T, 0, len(a)+len(b))

	i := 0
	j := 0
	for i < len(a) || j < len(b) {
		if i >= len(a) {
			out = out[:cap(out)]
			copy(out[i+j:], b[j:])
			break
		} else if j >= len(b) {
			out = out[:cap(out)]
			copy(out[i+j:], a[i:])
			break
		} else if cmp(a[i], b[j]) < 0 {
			out = append(out, a[i])
			i++
		} else {
			out = append(out, b[j])
			j++
		}
	}
	return out
}

func SortedMergeFunc[T any](in [][]T, cmp func(T, T) int) []T {
	if len(in) == 2 {
		return SortedMerge2Func(in[0], in[1], cmp)
	}
	n := 0
	for i := range in {
		n += len(in[i])
	}
	out := make([]T, 0, n)

	initial := make([]valueAndSource[T], 0, len(in))
	for i := range in {
		if len(in[i]) == 0 {
			continue
		}
		item := in[i][0]
		initial = append(initial, valueAndSource[T]{item, i})
	}

	h := heap.New(
		func(a, b valueAndSource[T]) int {
			return cmp(a.value, b.value)
		},
		func(a valueAndSource[T], i int) {},
		initial,
	)

	idxs := make([]int, len(in))
	for i := range idxs {
		idxs[i] = 1
	}

	for h.Len() > 0 {
		item := h.Pop()

		out = append(out, item.value)
		if idxs[item.source] < len(in[item.source]) {
			h.Push(valueAndSource[T]{
				value:  in[item.source][idxs[item.source]],
				source: item.source,
			})
			idxs[item.source]++
		}
	}
	return out
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
