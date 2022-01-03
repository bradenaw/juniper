package xtime

import (
	"context"
	"fmt"
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
