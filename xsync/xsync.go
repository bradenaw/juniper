package xsync

import (
	"context"
	"sync"

	"xsync/xatomic"
)

type ContextCond struct {
	ch xatomic.Value[chan struct{}]
	l sync.Locker
}

func NewContextCond(l sync.Locker) *ContextCond {
	c := ContextCond{
		l: l,
	}
	c.ch.Store(make(chan struct{}))
}

func (c *ContextCond) Broadcast() {
	ch := c.ch.Swap(make(chan struct{}))
	close(ch)
}

func (c *ContextCond) Signal() {
	select {
	case c.ch.Load() <- struct{}{}:
	default:
	}
}

func (c *ContextCond) Wait(ctx context.Context) error {
	ch := c.ch.Load()
	c.l.Unlock()
	defer c.l.Lock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
	}
}
