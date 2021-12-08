//go:build go1.18
// +build go1.18

package parallel

import (
	"context"
)

func Map[T any, U any](
	ctx context.Context,
	parallelism int,
	in []T,
	f func(ctx context.Context, in T) (U, error),
) ([]U, error) {
	out := make([]U, len(in))
	err := Do(ctx, parallelism, len(in), func(ctx context.Context, i int) error {
		var err error
		out[i], err = f(ctx, in[i])
		return err
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
