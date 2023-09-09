// Package xtime contains extensions to the standard library package time.
package xtime

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DeadlineTooSoonError struct {
	remaining time.Duration
	d         time.Duration
}

func (err DeadlineTooSoonError) Error() string {
	return fmt.Sprintf(
		"not enough time remaining in context: %s remaining for %s sleep",
		err.remaining,
		err.d,
	)
}

// SleepContext pauses the current goroutine for at least the duration d and returns nil, unless ctx
// expires in the mean time in which case it returns ctx.Err().
//
// A negative or zero duration causes SleepContext to return nil immediately.
//
// If there is less than d left until ctx's deadline, returns DeadlineTooSoonError immediately.
func SleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	deadline, ok := ctx.Deadline()
	if ok {
		remaining := time.Until(deadline)
		if remaining > d {
			return DeadlineTooSoonError{remaining: remaining, d: d}
		}
	}
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// A JitterTicker holds a channel that delivers "ticks" of a clock at intervals.
type JitterTicker struct {
	C <-chan time.Time

	c      chan time.Time
	m      sync.Mutex
	d      time.Duration
	gen    int
	jitter time.Duration
	timer  *time.Timer
}

// NewJitterTicker is similar to time.NewTicker, but jitters the ticks by the given amount. That is,
// each tick will be d+/-jitter apart.
//
// The duration d must be greater than zero and jitter must be less than d; if not, NewJitterTicker
// will panic.
func NewJitterTicker(d time.Duration, jitter time.Duration) *JitterTicker {
	if d <= 0 {
		panic("non-positive interval for NewJitterTicker")
	}
	if jitter >= d {
		panic("jitter greater than d")
	}

	c := make(chan time.Time, 1)
	t := &JitterTicker{
		C:      c,
		c:      c,
		d:      d,
		jitter: jitter,
	}
	t.m.Lock()
	t.schedule()
	t.m.Unlock()
	return t
}

func (t *JitterTicker) schedule() {
	if t.timer != nil {
		t.timer.Stop()
	}
	next := t.d + time.Duration(rand.Int63n(int64(t.jitter*2))) - (t.jitter)

	// To prevent a latent goroutine already spawned but not yet running the below function from
	// delivering a tick after Stop/Reset.
	t.gen++
	gen := t.gen

	t.timer = time.AfterFunc(next, func() {
		t.m.Lock()
		if t.gen == gen {
			select {
			case t.c <- time.Now():
			default:
			}
			t.schedule()
		}
		t.m.Unlock()
	})
}

// Reset stops the ticker and resets its period to be the specified duration and jitter. The next
// tick will arrive after the new period elapses.
//
// The duration d must be greater than zero and jitter must be less than d; if not, Reset will
// panic.
func (t *JitterTicker) Reset(d time.Duration, jitter time.Duration) {
	if d <= 0 {
		panic("non-positive interval for NewJitterTicker")
	}
	if jitter >= d {
		panic("jitter greater than d")
	}

	t.m.Lock()
	t.d = d
	t.jitter = jitter
	t.schedule()
	t.m.Unlock()
}

// Stop turns off the JitterTicker. After it returns, no more ticks will be sent. Stop does not
// close the channel, to prevent a concurrent goroutine reading from the channel from seeing an
// erroneous "tick".
func (t *JitterTicker) Stop() {
	t.m.Lock()
	t.timer.Stop()
	t.timer = nil
	t.m.Unlock()
}
