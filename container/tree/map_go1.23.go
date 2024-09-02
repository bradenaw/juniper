//go:build go1.23

package tree

import (
	"iter"

	"github.com/bradenaw/juniper/iterator"
)

func (m Map[K, V]) All() iter.Seq2[K, V] {
	return kvIterToSeq2(m.Iterate())
}

func (m Map[K, V]) Backward() iter.Seq2[K, V] {
	return kvIterToSeq2(m.RangeReverse(Unbounded[K](), Unbounded[K]()))
}

func (m Map[K, V]) Range2(lower Bound[K], upper Bound[K]) iter.Seq2[K, V] {
	return kvIterToSeq2(m.Range(lower, upper))
}

func (m Map[K, V]) RangeReverse2(lower Bound[K], upper Bound[K]) iter.Seq2[K, V] {
	return kvIterToSeq2(m.RangeReverse(lower, upper))
}

func kvIterToSeq2[K any, V any](it iterator.Iterator[KVPair[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for {
			kv, ok := it.Next()
			if !ok {
				break
			}
			if !yield(kv.Key, kv.Value) {
				break
			}
		}
	}
}
