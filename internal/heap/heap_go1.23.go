//go:build go1.23

package heap

import (
	"iter"
	"slices"
)

func (h *Heap[T]) All() iter.Seq[T] {
	return slices.Values(h.a)
}
