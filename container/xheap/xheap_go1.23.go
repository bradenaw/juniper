//go:build go1.23

package xheap

import (
	"iter"
)

func (h PriorityQueue[K, P]) All() iter.Seq2[K, P] {
	return func(yield func(K, P) bool) {
		for item := range h.inner.All() {
			if !yield(item.K, item.P) {
				return
			}
		}
	}
}
