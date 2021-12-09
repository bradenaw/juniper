// Package xsync contains extensions to the standard library package sync.
package xsync

import (
	"context"
	"sync"
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
