// Package xsync contains extensions to the standard library package sync.
package xsync

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// ContextCond is equivalent to sync.Cond, except its Wait function accepts a context.Context.
//
// ContextConds should not be copied after first use.
type ContextCond struct {
	m  sync.RWMutex
	ch chan struct{}
	L  sync.Locker
}

// NewContextCond returns a new ContextCond with l as its Locker.
func NewContextCond(l sync.Locker) *ContextCond {
	return &ContextCond{
		L:  l,
		ch: make(chan struct{}),
	}
}

// Broadcast wakes all goroutines blocked in Wait(), if there are any.
//
// It is allowed but not required for the caller to hold c.L during the call.
func (c *ContextCond) Broadcast() {
	c.m.Lock()
	close(c.ch)
	c.ch = make(chan struct{})
	c.m.Unlock()
}

// Signal wakes one goroutine blocked in Wait(), if there is any. No guarantee is made as to which
// goroutine will wake.
//
// It is allowed but not required for the caller to hold c.L during the call.
func (c *ContextCond) Signal() {
	c.m.RLock()
	select {
	case c.ch <- struct{}{}:
	default:
	}
	c.m.RUnlock()
}

// Wait is equivalent to sync.Cond.Wait, except it accepts a context.Context. If the context expires
// before this goroutine is woken by Broadcast or Signal, it returns ctx.Err() immediately. If an
// error is returned, does not reaquire c.L before returning.
func (c *ContextCond) Wait(ctx context.Context) error {
	c.m.RLock()
	ch := c.ch
	c.m.RUnlock()
	c.L.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		c.L.Lock()
	}
	return nil
}

// Group manages a group of goroutines.
type Group struct {
	baseCtx context.Context
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewGroup returns a Group ready for use. The context passed to any of the f functions will be a
// descendant of ctx.
func NewGroup(ctx context.Context) *Group {
	bgCtx, cancel := context.WithCancel(ctx)
	return &Group{
		baseCtx: ctx,
		ctx:     bgCtx,
		cancel:  cancel,
	}
}

// Once calls f once from another goroutine.
func (g *Group) Once(f func(ctx context.Context)) {
	g.wg.Add(1)
	go func() {
		f(g.ctx)
		g.wg.Done()
	}()
}

// returns a random duration in [d - jitter, d + jitter]
func jitterDuration(d time.Duration, jitter time.Duration) time.Duration {
	return d + time.Duration(float64(jitter)*((rand.Float64()*2)-1))
}

// Periodic spawns a goroutine that calls f once per interval +/- jitter.
func (g *Group) Periodic(
	interval time.Duration,
	jitter time.Duration,
	f func(ctx context.Context),
) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		t := time.NewTimer(jitterDuration(interval, jitter))
		defer t.Stop()
		for {
			if g.ctx.Err() != nil {
				return
			}
			select {
			case <-g.ctx.Done():
				return
			case <-t.C:
			}
			t.Reset(jitterDuration(interval, jitter))
			f(g.ctx)
		}
	}()
}

// Trigger spawns a goroutine which calls f whenever the returned function is called. If f is
// already running when triggered, f will run again immediately when it finishes.
func (g *Group) Trigger(f func(ctx context.Context)) func() {
	c := make(chan struct{}, 1)
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		for {
			if g.ctx.Err() != nil {
				return
			}
			select {
			case <-g.ctx.Done():
				return
			case <-c:
			}
			f(g.ctx)
		}
	}()

	return func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

// PeriodicOrTrigger spawns a goroutine which calls f whenever the returned function is called.  If
// f is already running when triggered, f will run again immediately when it finishes. Also calls f
// when it has been interval+/-jitter since the last trigger.
func (g *Group) PeriodicOrTrigger(
	interval time.Duration,
	jitter time.Duration,
	f func(ctx context.Context),
) func() {
	c := make(chan struct{}, 1)
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		t := time.NewTimer(jitterDuration(interval, jitter))
		defer t.Stop()
		for {
			if g.ctx.Err() != nil {
				return
			}
			select {
			case <-g.ctx.Done():
				return
			case <-t.C:
				t.Reset(jitterDuration(interval, jitter))
			case <-c:
				if !t.Stop() {
					<-t.C
				}
				t.Reset(jitterDuration(interval, jitter))
			}
			f(g.ctx)
		}
	}()

	return func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

// Stop cancels the context passed to spawned goroutines.
func (g *Group) Stop() {
	g.cancel()
}

// Wait cancels the context passed to any of the spawned goroutines and waits for all spawned
// goroutines to exit.
//
// It is not safe to call Wait concurrently with any other method on g.
func (g *Group) Wait() {
	g.cancel()
	g.wg.Wait()
	g.ctx, g.cancel = context.WithCancel(g.baseCtx)
}

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
	if !ok {
		var zero V
		return zero, false
	}
	return value_.(V), ok
}
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	value_, ok := m.m.LoadAndDelete(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value_.(V), ok
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
