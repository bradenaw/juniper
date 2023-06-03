package xsync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bradenaw/juniper/xtime"
)

func ExampleLazy() {
	var (
		expensive = Lazy(func() string {
			fmt.Println("doing expensive init")
			return "foo"
		})
	)

	fmt.Println(expensive())
	fmt.Println(expensive())

	// Output:
	// doing expensive init
	// foo
	// foo
}

func TestGroup(t *testing.T) {
	g := NewGroup(context.Background())

	dos := make(chan struct{}, 100)
	g.Do(func(ctx context.Context) {
		for {
			err := xtime.SleepContext(ctx, 50*time.Millisecond)
			if err != nil {
				return
			}

			select {
			case dos <- struct{}{}:
			default:
			}
		}
	})

	periodics := make(chan struct{}, 100)
	g.Periodic(35*time.Millisecond, 0 /*jitter*/, func(ctx context.Context) {
		select {
		case periodics <- struct{}{}:
		default:
		}
	})

	periodicOrTriggers := make(chan struct{}, 100)
	periodicOrTrigger := g.PeriodicOrTrigger(75*time.Millisecond, 0 /*jitter*/, func(ctx context.Context) {
		select {
		case periodicOrTriggers <- struct{}{}:
		default:
		}
	})

	triggers := make(chan struct{}, 100)
	trigger := g.Trigger(func(ctx context.Context) {
		select {
		case triggers <- struct{}{}:
		default:
		}
	})

	trigger()
	periodicOrTrigger()
	time.Sleep(200 * time.Millisecond)
	trigger()

	<-dos
	<-dos
	<-dos
	<-dos
	<-periodics
	<-periodics
	<-periodics
	<-periodics
	<-periodics
	<-periodicOrTriggers
	<-periodicOrTriggers
	<-periodicOrTriggers
	<-triggers
	<-triggers

	g.StopAndWait()

	g.Do(func(ctx context.Context) {
		panic("this will never spawn because StopAndWait was already called")
	})

	// Jank, but just in case we'd be safe from the above panic just because the test is over.
	time.Sleep(200 * time.Millisecond)
}
