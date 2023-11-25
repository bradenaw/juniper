//go:build go1.19

package xsync

import (
	"sync/atomic"
)

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
