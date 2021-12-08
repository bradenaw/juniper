package xsort

import (
	"constraints"
	"sort"
)

type Less[T any] func(a, b T) bool

func OrderedLess[T constraints.Ordered](a, b T) bool {
	return a < b
}

var _ Less[int] = OrderedLess[int]

func Slice[T any](x []T, less Less[T]) {
	sort.Slice(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

func SliceStable[T any](x []T, less Less[T]) {
	sort.SliceStable(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}

func SliceIsSorted[T any](x []T, less Less[T]) bool {
	return sort.SliceIsSorted(x, func(i, j int) bool {
		return less(x[i], x[j])
	})
}
