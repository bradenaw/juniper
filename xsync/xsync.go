package xsync

import (
	"context"
	"sync"
)

type ContextCond struct {
	m  sync.RWMutex
	ch chan struct{}
	l  sync.Locker
}

func NewContextCond(l sync.Locker) *ContextCond {
	return &ContextCond{
		l:  l,
		ch: make(chan struct{}),
	}
}

func (c *ContextCond) Broadcast() {
	c.m.Lock()
	close(c.ch)
	c.ch = make(chan struct{})
	c.m.Unlock()
}

func (c *ContextCond) Signal() {
	c.m.RLock()
	select {
	case c.ch <- struct{}{}:
	default:
	}
	c.m.RUnlock()
}

func (c *ContextCond) Wait(ctx context.Context) error {
	c.m.RLock()
	ch := c.ch
	c.m.RUnlock()
	c.l.Unlock()
	defer c.l.Lock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
	}
	return nil
}
