//go:build go1.18

package chans

import "context"

// SendOrExpire sends item on channel c and returns nil, unless ctx expires in which case it returns
// ctx.Err().
func SendOrExpire[T any](ctx context.Context, c chan<- T, item T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c <- item:
		return nil
	}
}

// RecvOrExpire attempts to receive from channel c. If c is closed before or during, returns (_,
// false, nil). If ctx expires before or during, returns (_, _, ctx.Err()).
func RecvOrExpire[T any](ctx context.Context, c <-chan T) (T, bool, error) {
	select {
	case <-ctx.Done():
		var zero T
		return zero, false, ctx.Err()
	case item, ok := <-c:
		return item, ok, nil
	}
}
