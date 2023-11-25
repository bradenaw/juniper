// Package xsync contains extensions to the standard library package sync.
package xsync

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
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
	ctx    context.Context
	cancel context.CancelFunc
	// held in R when spawning to check if ctx is already cancelled and in W when cancelling ctx to
	// make sure we never cause wg to go 0->1 while inside Wait()
	m  sync.RWMutex
	wg sync.WaitGroup
}

// NewGroup returns a Group ready for use. The context passed to any of the f functions will be a
// descendant of ctx.
func NewGroup(ctx context.Context) *Group {
	bgCtx, cancel := context.WithCancel(ctx)
	return &Group{
		ctx:    bgCtx,
		cancel: cancel,
	}
}

// helper even though it's exactly g.Do so that the goroutine stack for a spawned function doesn't
// confusingly show all of them as created by Do.
func (g *Group) spawn(f func()) {
	g.m.RLock()
	if g.ctx.Err() != nil {
		g.m.RUnlock()
		return
	}
	g.wg.Add(1)
	g.m.RUnlock()

	go func() {
		f()
		g.wg.Done()
	}()
}

// Do calls f once from another goroutine.
func (g *Group) Do(f func(ctx context.Context)) {
	g.spawn(func() { f(g.ctx) })
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
	g.spawn(func() {
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
	})
}

// Trigger spawns a goroutine which calls f whenever the returned function is called. If f is
// already running when triggered, f will run again immediately when it finishes.
func (g *Group) Trigger(f func(ctx context.Context)) func() {
	c := make(chan struct{}, 1)

	g.spawn(func() {
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
	})

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
	g.spawn(func() {
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
	})

	return func() {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

// Stop cancels the context passed to spawned goroutines. After the group is stopped, no more
// goroutines will be spawned.
func (g *Group) Stop() {
	g.m.Lock()
	g.cancel()
	g.m.Unlock()
}

// StopAndWait cancels the context passed to any of the spawned goroutines and waits for all spawned
// goroutines to exit. After the group is stopped, no more goroutines will be spawned.
func (g *Group) StopAndWait() {
	g.Stop()
	g.wg.Wait()
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

// Watchable contains a value. It is similar to an atomic.Pointer[T] but allows notifying callers
// that a new value has been set.
type Watchable[T any] struct {
	p atomic.Pointer[watchableInner[T]]
}

type watchableInner[T any] struct {
	t T
	c chan struct{}
}

// Set sets the value in w and notifies callers of Value() that there is a new value.
func (w *Watchable[T]) Set(t T) {
	newInner := &watchableInner[T]{
		t: t,
		c: make(chan struct{}),
	}
	oldInner := w.p.Swap(newInner)
	if oldInner != nil {
		close(oldInner.c)
	}
}

// Value returns the current value inside w and a channel that will be closed when w is Set() to a
// newer value than the returned one.
//
// If called before the first Set(), returns the zero value of T.
//
// Normal usage has an observer waiting for new values in a loop:
//
//	for {
//		v, changed := w.Value()
//
//		// do something with v
//
//		<-changed
//	}
//
// Note that the value in w may have changed multiple times between successive calls to Value(),
// Value() only ever returns the last-set value. This is by design so that slow observers cannot
// block Set(), unlike sending values on a channel.
func (w *Watchable[T]) Value() (T, chan struct{}) {
	inner := w.p.Load()
	if inner == nil {
		// There's no inner, meaning w has not been Set() yet. Try filling it with an empty inner,
		// so that we have a channel to listen on.
		c := make(chan struct{})
		emptyInner := &watchableInner[T]{
			c: c,
		}
		// CompareAndSwap so we don't accidentally smash a real value that got put between our Load
		// and here.
		if w.p.CompareAndSwap(nil, emptyInner) {
			var zero T
			return zero, c
		}
		// If we fell through to here somebody Set() while we were trying to do this, so there's
		// definitely an inner now.
		inner = w.p.Load()
	}
	return inner.t, inner.c
}
