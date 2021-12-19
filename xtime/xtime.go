package xtime

import (
	"context"
	"time"
)

// SleepContext pauses the current goroutine for at least the duration d and returns nil, unless ctx
// expires in the mean time in which case it returns ctx.Err().
//
// A negative or zero duration causes SleepContext to return immediately.
func SleepContext(ctx context.Context, d time.Duration) error {
	if d < 0 {
		return nil
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
