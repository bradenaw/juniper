//go:build go1.18

package xsync

import "sync"

// Lazy makes a lazily-initialized value. On first access, it uses f to create the value. Later
// accesses all receive the same value.
func Lazy[T any](f func() T) func() T {
	var once sync.Once
	var val T
	return func() T {
		once.Do(func() {
			val = f()
		})
		return val
	}
}

// Map is a typesafe wrapper over sync.Map.
type Map[K comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	value_, ok := m.m.Load(key)
	return value_.(V), ok
}
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	value_, loaded := m.m.LoadAndDelete(key)
	return value_.(V), loaded
}
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	actual_, loaded := m.m.LoadOrStore(key, value)
	return actual_.(V), loaded
}
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Pool is a typesafe wrapper over sync.Pool.
type Pool[T any] struct {
	sync.Pool
}

func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.Pool.Put(x)
}
