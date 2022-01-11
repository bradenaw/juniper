//go:build go1.18

package xsync

import (
	"context"
	"sync"
)

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
	p sync.Pool
}

func NewPool[T any](new_ func() T) Pool[T] {
	return Pool[T]{
		p: sync.Pool{
			New: func() interface{} {
				return new_()
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.p.Put(x)
}

// Future can be filled with a value exactly once. Many goroutines can concurrently wait for it to
// be filled. After filling, Wait() immediately returns the value it was filled with.
//
// Futures must be created by NewFuture and should not be copied after first use.
type Future[T any] struct {
	c chan struct{}
	x T
}

// NewFuture returns a ready-to-use Future.
func NewFuture[T any]() *Future[T] {
	return &Future[T]{
		c: make(chan struct{}),
	}
}

// Fill fills f with value x. All active calls to Wait return x, and all future calls to Wait return
// x immediately.
//
// Panics if f has already been filled.
func (f *Future[T]) Fill(x T) {
	f.x = x
	close(f.c)
}

// Wait waits for f to be filled with a value and returns it. Returns immediately if f is already
// filled.
func (f *Future[T]) Wait() T {
	<-f.c
	return f.x
}

// Wait waits for f to be filled with a value and returns it, or returns ctx.Err() if ctx expires
// before this happens. Returns immediately if f is already filled.
func (f *Future[T]) WaitContext(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case <-f.c:
	}
	return f.x, nil
}
