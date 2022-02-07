//go:build go1.18

package chans

import (
	"context"
	"reflect"

	"github.com/bradenaw/juniper/slices"
)

// SendContext sends item on channel c and returns nil, unless ctx expires in which case it returns
// ctx.Err().
func SendContext[T any](ctx context.Context, c chan<- T, item T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c <- item:
		return nil
	}
}

// RecvContext attempts to receive from channel c. If c is closed before or during, returns (_,
// false, nil). If ctx expires before or during, returns (_, _, ctx.Err()).
func RecvContext[T any](ctx context.Context, c <-chan T) (T, bool, error) {
	select {
	case <-ctx.Done():
		var zero T
		return zero, false, ctx.Err()
	case item, ok := <-c:
		return item, ok, nil
	}
}

// Merge sends all values from all in channels to out.
//
// Merge blocks until all ins have closed and all values have been sent. It does not close out.
func Merge[T any](out chan<- T, in ...<-chan T) {
	if len(in) == 1 {
		for item := range in[0] {
			out <- item
		}
	} else if len(in) == 2 {
		merge2(out, in[0], in[1])
		return
	} else if len(in) == 3 {
		merge3(out, in[0], in[1], in[2])
		return
	}

	selectCases := slices.Map(in, func(x <-chan T) reflect.SelectCase {
		return reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(x),
		}
	})
	for {
		if len(selectCases) == 0 {
			return
		}
		chosen, item, ok := reflect.Select(selectCases)
		if ok {
			out <- item.Interface().(T)
		} else {
			selectCases = slices.RemoveUnordered(selectCases, chosen, 1)
		}
	}
}

// Merge special-case with no reflection.
func merge2[T any](out chan<- T, in0, in1 <-chan T) {
	nDone := 0
	for {
		select {
		case item, ok := <-in0:
			if ok {
				out <- item
			} else {
				in0 = nil
				nDone++
				if nDone == 2 {
					return
				}
			}
		case item, ok := <-in1:
			if ok {
				out <- item
			} else {
				in1 = nil
				nDone++
				if nDone == 2 {
					return
				}
			}
		}
	}
}

// Merge special-case with no reflection.
func merge3[T any](out chan<- T, in0, in1, in2 <-chan T) {
	nDone := 0
	for {
		select {
		case item, ok := <-in0:
			if ok {
				out <- item
			} else {
				in0 = nil
				nDone++
				if nDone == 3 {
					return
				}
			}
		case item, ok := <-in1:
			if ok {
				out <- item
			} else {
				in1 = nil
				nDone++
				if nDone == 3 {
					return
				}
			}
		case item, ok := <-in2:
			if ok {
				out <- item
			} else {
				in2 = nil
				nDone++
				if nDone == 3 {
					return
				}
			}
		}
	}
}

// Replicate sends all values sent to src to every channel in dsts.
//
// Replicate blocks until src is closed and all values have been sent to all dsts. It does not close
// dsts.
func Replicate[T any](src <-chan T, dsts ...chan<- T) {
	for item := range src {
		for _, dst := range dsts {
			dst <- item
		}
	}
}
