//go:build go1.18

// package orderedhashmap contains a simple and very inefficient ordered map using the map builtin
// for comparing against other ordered containers in tests.
package orderedhashmap

import (
	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/xsort"
)

type KVPair[K any, V any] struct {
	K K
	V V
}

type Map[K comparable, V any] struct {
	less xsort.Less[K]
	m    map[K]V
}

func NewMap[K comparable, V any](less xsort.Less[K]) Map[K, V] {
	return Map[K, V]{
		less: less,
		m:    make(map[K]V),
	}
}

func (m Map[K, V]) Len() int {
	return len(m.m)
}

func (m Map[K, V]) Put(k K, v V) {
	m.m[k] = v
}

func (m Map[K, V]) Delete(k K) {
	delete(m.m, k)
}

func (m Map[K, V]) Get(k K) V {
	return m.m[k]
}

func (m Map[K, V]) Contains(k K) bool {
	_, ok := m.m[k]
	return ok
}

func (m Map[K, V]) First() (K, V) {
	first := true
	var min K
	for k := range m.m {
		if first || m.less(k, min) {
			min = k
			first = false
		}
	}
	return min, m.m[min]
}

func (m Map[K, V]) Last() (K, V) {
	first := true
	var max K
	for k := range m.m {
		if first || m.less(max, k) {
			max = k
			first = false
		}
	}
	return max, m.m[max]
}

func (m Map[K, V]) Iterate() iterator.Iterator[KVPair[K, V]] {
	return m.Cursor().Forward()
}

func (m Map[K, V]) Cursor() *MapCursor[K, V] {
	c := &MapCursor[K, V]{
		m: m,
	}
	c.SeekFirst()
	return c
}

func (m Map[K, V]) lastLess(k K) (K, bool) {
	first := true
	var out K
	for existingK := range m.m {
		if xsort.GreaterOrEqual(m.less, existingK, k) {
			continue
		}
		if first || m.less(out, existingK) {
			out = existingK
			first = false
		}
	}
	return out, !first
}

func (m Map[K, V]) firstGreater(k K) (K, bool) {
	first := true
	var out K
	for existingK := range m.m {
		if xsort.LessOrEqual(m.less, existingK, k) {
			continue
		}
		if first || m.less(existingK, out) {
			out = existingK
			first = false
		}
	}
	return out, !first
}

type MapCursor[K comparable, V any] struct {
	m  Map[K, V]
	ok bool
	k  K
}

func (c *MapCursor[K, V]) SeekFirst() {
	c.k, _ = c.m.First()
	c.ok = len(c.m.m) > 0
}

func (c *MapCursor[K, V]) SeekLast() {
	c.k, _ = c.m.Last()
	c.ok = len(c.m.m) > 0
}

func (c *MapCursor[K, V]) set(k K) {
	c.k = k
	c.ok = true
}

func (c *MapCursor[K, V]) SeekLastLess(k K) {
	k, ok := c.m.lastLess(k)
	c.ok = ok
	if ok {
		c.set(k)
	}
}

func (c *MapCursor[K, V]) SeekLastLessOrEqual(k K) {
	if c.m.Contains(k) {
		c.set(k)
		return
	}
	k, ok := c.m.lastLess(k)
	c.ok = ok
	if ok {
		c.set(k)
	}
}

func (c *MapCursor[K, V]) SeekFirstGreaterOrEqual(k K) {
	if c.m.Contains(k) {
		c.set(k)
		return
	}
	k, ok := c.m.firstGreater(k)
	c.ok = ok
	if ok {
		c.set(k)
	}
}

func (c *MapCursor[K, V]) SeekFirstGreater(k K) {
	k, ok := c.m.firstGreater(k)
	c.ok = ok
	if ok {
		c.set(k)
	}
}

func (c *MapCursor[K, V]) Next() {
	if !c.ok {
		return
	}
	k, ok := c.m.firstGreater(c.k)
	c.ok = ok
	if ok {
		c.set(k)
	}
}

func (c *MapCursor[K, V]) Prev() {
	if !c.ok {
		return
	}
	k, ok := c.m.lastLess(c.k)
	c.ok = ok
	if ok {
		c.set(k)
	}
}

func (c *MapCursor[K, V]) deleted() bool {
	return !c.m.Contains(c.k)
}

func (c *MapCursor[K, V]) Ok() bool { return c.ok }

func (c *MapCursor[K, V]) Key() K { return c.k }

func (c *MapCursor[K, V]) Value() V { return c.m.m[c.k] }

type forwardIterator[K comparable, V any] struct {
	c MapCursor[K, V]
}

func (iter *forwardIterator[K, V]) Next() (KVPair[K, V], bool) {
	if !iter.c.Ok() {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	v := iter.c.Value()
	iter.c.Next()
	return KVPair[K, V]{k, v}, true
}

func (c *MapCursor[K, V]) Forward() iterator.Iterator[KVPair[K, V]] {
	c2 := *c
	if c2.Ok() && c2.deleted() {
		c2.SeekFirstGreater(c2.k)
	}
	return &forwardIterator[K, V]{c: c2}
}

type backwardIterator[K comparable, V any] struct {
	c MapCursor[K, V]
}

func (iter *backwardIterator[K, V]) Next() (KVPair[K, V], bool) {
	if !iter.c.Ok() {
		var zero KVPair[K, V]
		return zero, false
	}
	k := iter.c.Key()
	v := iter.c.Value()
	iter.c.Prev()
	return KVPair[K, V]{k, v}, true
}

func (c *MapCursor[K, V]) Backward() iterator.Iterator[KVPair[K, V]] {
	c2 := *c
	if c2.Ok() && c2.deleted() {
		c2.SeekLastLess(c.k)
	}
	return &backwardIterator[K, V]{c: c2}
}
