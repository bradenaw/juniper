package xtime

import (
	"testing"
	"time"
)

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
