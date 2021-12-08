package xsort

import (
	"constraints"
	"sort"

	"github.com/bradenaw/xstd/internal/heap"
	"github.com/bradenaw/xstd/slices"
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

type valueAndSource[T any] struct {
	value  T
	source int
}

func Merge[T any](less Less[T], out []T, in ...[]T) []T {
	initial := make([]valueAndSource[T], 0, len(in))
	n := 0
	for j := range in {
		n += len(in[j])
		if len(in[j]) > 0 {
			initial = append(initial, valueAndSource[T]{in[j][0], j})
			in[j] = in[j][1:]
		}
	}
	h := heap.New(
		func(a, b valueAndSource[T]) bool {
			return less(a.value, b.value)
		},
		func(a valueAndSource[T], i int) {},
		initial,
	)
	out = slices.Grow(out[:0], n)
	for {
		item, ok := h.Pop()
		if !ok {
			break
		}
		out = append(out, item.value)
		if len(in[item.source]) > 0 {
			h.Push(valueAndSource[T]{in[item.source][0], item.source})
			in[item.source] = in[item.source][1:]
		}
	}
	return out
}
