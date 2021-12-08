package parallel

import (
	"context"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

func Do(
	ctx context.Context,
	parallelism int,
	n int,
	f func(ctx context.Context, i int) error,
) error {
	if parallelism == 1 {
		for i := 0; i < n; i++ {
			err := f(ctx, i)
			if err != nil {
				return err
			}
		}
		return nil
	}

	x := int32(-1)
	eg, ctx := errgroup.WithContext(ctx)
	for j := 0; j < parallelism; j++ {
		eg.Go(func() error {
			for {
				i := int(atomic.AddInt32(&x, 1))
				if i >= n {
					return nil
				}

				err := f(ctx, i)
				if err != nil {
					return err
				}
			}
		})
	}
	return eg.Wait()
}
