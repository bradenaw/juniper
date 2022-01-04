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
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewGroup returns a Group ready for use.
func NewGroup() *Group {
	ctx, cancel := context.WithCancel(context.Background())
	return &Group{
		ctx:    ctx,
		cancel: cancel,
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
//
// It is not safe to call Stop concurrently with any other method on g.
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
	g.ctx, g.cancel = context.WithCancel(context.Background())
}
