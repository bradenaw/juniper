package xtime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSleepContext(t *testing.T) {
	t.Run("sufficient deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := SleepContext(ctx, time.Millisecond); err != nil {
			t.Fatalf("SleepContext() = %v, want nil", err)
		}
	})

	t.Run("insufficient deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := SleepContext(ctx, 2*time.Second)
		var deadlineErr DeadlineTooSoonError
		if !errors.As(err, &deadlineErr) {
			t.Fatalf("SleepContext() = %v, want DeadlineTooSoonError", err)
		}
		if err := ctx.Err(); err != nil {
			t.Fatalf("context expired before SleepContext returned: %v", err)
		}
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if err := SleepContext(ctx, time.Hour); !errors.Is(err, context.Canceled) {
			t.Fatalf("SleepContext() = %v, want context.Canceled", err)
		}
	})

	for _, d := range []time.Duration{0, -time.Nanosecond} {
		t.Run(d.String(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			if err := SleepContext(ctx, d); err != nil {
				t.Fatalf("SleepContext() = %v, want nil", err)
			}
		})
	}
}

func TestJitterTicker(t *testing.T) {
	d := 5 * time.Millisecond
	jitter := 2 * time.Millisecond
	ticker := NewJitterTicker(d, jitter)

	last := time.Now()
	check := func() {
		now := time.Now()
		elapsed := now.Sub(last)
		minTick := d - jitter
		// Add a little extra slack because of scheduling.
		maxTick := d + jitter + 3*time.Millisecond

		if elapsed < minTick {
			t.Fatalf("tick was %s, expected in [%s, %s]", elapsed, minTick, maxTick)
		}
		if elapsed > maxTick {
			t.Fatalf("tick was %s, expected in [%s, %s]", elapsed, minTick, maxTick)
		}

		last = now
	}

	for i := 0; i < 50; i++ {
		<-ticker.C
		check()
	}

	d = 10 * time.Millisecond
	jitter = 8 * time.Millisecond
	ticker.Reset(d, jitter)
	for i := 0; i < 20; i++ {
		<-ticker.C
		check()
	}
}
