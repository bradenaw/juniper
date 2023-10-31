//go:build go1.21

package xsync

import (
	"sync"
)

// Lazy makes a lazily-initialized value. On first access, it uses f to create the value. Later
// accesses all receive the same value.
//
// Deprecated: sync.OnceValue is in the standard library as of Go 1.21.
func Lazy[T any](f func() T) func() T {
	return sync.OnceValue(f)
}

func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}
func (m *Map[K, V]) CompareAndSwap(key K, old V, new V) (deleted bool) {
	return m.m.CompareAndSwap(key, old, new)
}
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	previousUntyped, loaded := m.m.Swap(key, value)
	return previousUntyped.(V), loaded
}
