//go:build go1.19

package xsync

import (
	"fmt"
	"time"
)

func ExampleWatchable() {
	start := time.Now()

	var w Watchable[int]
	w.Set(0)
	go func() {
		for i := 1; i < 20; i++ {
			w.Set(i)
			fmt.Printf("set %d at %s\n", i, time.Since(start).Round(time.Millisecond))
			time.Sleep(5 * time.Millisecond)
		}
	}()

	for {
		v, changed := w.Value()
		if v == 19 {
			return
		}

		fmt.Printf("observed %d at %s\n", v, time.Since(start).Round(time.Millisecond))

		// Sleep for longer between iterations to show that we don't slow down the setter.
		time.Sleep(17 * time.Millisecond)

		<-changed
	}
}
