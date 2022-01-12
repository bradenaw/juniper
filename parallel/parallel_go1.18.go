//go:build go1.18
// +build go1.18

package parallel

import (
	"context"
	"runtime"

	"golang.org/x/sync/errgroup"

	"github.com/bradenaw/juniper/iterator"
	"github.com/bradenaw/juniper/slices"
	"github.com/bradenaw/juniper/stream"
	"github.com/bradenaw/juniper/xsort"
)

// Map uses parallelism goroutines to call f once for each element of in. out[i] is the
// result of f for in[i].
//
// If parallelism <= 0, uses GOMAXPROCS instead.
func Map[T any, U any](
	parallelism int,
	in []T,
	f func(in T) U,
) []U {
	out := make([]U, len(in))
	Do(parallelism, len(in), func(i int) {
		out[i] = f(in[i])
	})
	return out
}

// MapContext uses parallelism goroutines to call f once for each element of in. out[i] is the
// result of f for in[i].
//
// If any call to f returns an error the context passed to invocations of f is cancelled, no further
// calls to f are made, and Map returns the first error encountered.
//
// If parallelism <= 0, uses GOMAXPROCS instead.
func MapContext[T any, U any](
	ctx context.Context,
	parallelism int,
	in []T,
	f func(ctx context.Context, in T) (U, error),
) ([]U, error) {
	out := make([]U, len(in))
	err := DoContext(ctx, parallelism, len(in), func(ctx context.Context, i int) error {
		var err error
		out[i], err = f(ctx, in[i])
		return err
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapIterator uses parallelism goroutines to call f once for each element yielded by iter. The
// returned iterator returns these results in the same order that iter yielded them in.
//
// This iterator, in contrast with most, must be consumed completely or it will leak the goroutines.
//
// If parallelism <= 0, uses GOMAXPROCS instead.
//
// bufferSize is the size of the work buffer for each goroutine. A larger buffer uses more memory
// but gives better throughput in the face of larger variance in the processing time for f.
func MapIterator[T any, U any](
	iter iterator.Iterator[T],
	parallelism int,
	bufferSize int,
	f func(T) U,
) iterator.Iterator[U] {
	if parallelism <= 0 {
		parallelism = runtime.GOMAXPROCS(-1)
	}

	in := make(chan valueAndIndex[T])

	go func() {
		i := 0
		for {
			item, ok := iter.Next()
			if !ok {
				break
			}

			in <- valueAndIndex[T]{
				value: item,
				idx:   i,
			}
			i++
		}
		close(in)
	}()

	outs := make([]chan valueAndIndex[U], parallelism)
	for i := 0; i < parallelism; i++ {
		i := i
		outs[i] = make(chan valueAndIndex[U], bufferSize)
		go func() {
			for item := range in {
				u := f(item.value)
				outs[i] <- valueAndIndex[U]{value: u, idx: item.idx}
			}
			close(outs[i])
		}()
	}

	return iterator.Map(
		xsort.Merge(
			func(a, b valueAndIndex[U]) bool {
				return a.idx < b.idx
			},
			slices.Map(outs, func(c chan valueAndIndex[U]) iterator.Iterator[valueAndIndex[U]] {
				return iterator.Chan(c)
			})...,
		),
		func(x valueAndIndex[U]) U {
			return x.value
		},
	)
}

type valueAndIndex[T any] struct {
	value T
	idx   int
}

// MapStream uses parallelism goroutines to call f once for each element yielded by s. The returned
// stream returns these results in the same order that s yielded them in.
//
// If any call to f returns an error the context passed to invocations of f is cancelled, no further
// calls to f are made, and the returned stream's Next returns the first error encountered.
//
// If parallelism <= 0, uses GOMAXPROCS instead.
//
// bufferSize is the size of the work buffer for each goroutine. A larger buffer uses more memory
// but gives better throughput in the face of larger variance in the processing time for f.
func MapStream[T any, U any](
	s stream.Stream[T],
	parallelism int,
	bufferSize int,
	f func(context.Context, T) (U, error),
) stream.Stream[U] {
	if parallelism <= 0 {
		parallelism = runtime.GOMAXPROCS(-1)
	}

	in := make(chan valueAndIndex[T])

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		defer s.Close()
		defer close(in)
		i := 0
		for {
			item, err := s.Next(ctx)
			if err == stream.End {
				break
			} else if err != nil {
				return err
			}

			select {
			case in <- valueAndIndex[T]{
				value: item,
				idx:   i,
			}:
			case <-ctx.Done():
				return ctx.Err()
			}
			i++
		}
		return nil
	})

	outs := make([]chan valueAndIndex[U], parallelism)
	for i := 0; i < parallelism; i++ {
		i := i
		outs[i] = make(chan valueAndIndex[U], bufferSize)
		eg.Go(func() error {
			defer close(outs[i])
			for item := range in {
				u, err := f(ctx, item.value)
				if err != nil {
					return err
				}
				select {
				case outs[i] <- valueAndIndex[U]{value: u, idx: item.idx}:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	sender, receiver := stream.Pipe[U](0)

	go func() {
		iter := iterator.Map(
			xsort.Merge(
				func(a, b valueAndIndex[U]) bool {
					return a.idx < b.idx
				},
				slices.Map(outs, func(c chan valueAndIndex[U]) iterator.Iterator[valueAndIndex[U]] {
					return iterator.Chan(c)
				})...,
			),
			func(x valueAndIndex[U]) U {
				return x.value
			},
		)
		for {
			item, ok := iter.Next()
			if !ok {
				break
			}
			err := sender.Send(ctx, item)
			if err != nil {
				// could be either the receiver hung up or one of the above goroutines ran into a
				// problem (which would've cancelled ctx already), either way it's time to bail out
				cancel()
				break
			}
		}
		sender.Close(eg.Wait())
	}()

	return receiver
}
