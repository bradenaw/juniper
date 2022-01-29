# parallel
--
    import "."

Package parallel provides primitives for running tasks in parallel.

## Usage

#### func  Do

```go
func Do(
	parallelism int,
	n int,
	f func(i int),
)
```
Do calls f from parallelism goroutines n times, providing each invocation a
unique i in [0, n).

If parallelism <= 0, uses GOMAXPROCS instead.

#### func  DoContext

```go
func DoContext(
	ctx context.Context,
	parallelism int,
	n int,
	f func(ctx context.Context, i int) error,
) error
```
DoContext calls f from parallelism goroutines n times, providing each invocation
a unique i in [0, n).

If any call to f returns an error the context passed to invocations of f is
cancelled, no further calls to f are made, and Do returns the first error
encountered.

If parallelism <= 0, uses GOMAXPROCS instead.

#### func  Map

```go
func Map[T any, U any](
	parallelism int,
	in []T,
	f func(in T) U,
) []U
```
Map uses parallelism goroutines to call f once for each element of in. out[i] is
the result of f for in[i].

If parallelism <= 0, uses GOMAXPROCS instead.

#### func  MapContext

```go
func MapContext[T any, U any](
	ctx context.Context,
	parallelism int,
	in []T,
	f func(ctx context.Context, in T) (U, error),
) ([]U, error)
```
MapContext uses parallelism goroutines to call f once for each element of in.
out[i] is the result of f for in[i].

If any call to f returns an error the context passed to invocations of f is
cancelled, no further calls to f are made, and Map returns the first error
encountered.

If parallelism <= 0, uses GOMAXPROCS instead.

#### func  MapIterator

```go
func MapIterator[T any, U any](
	iter iterator.Iterator[T],
	parallelism int,
	bufferSize int,
	f func(T) U,
) iterator.Iterator[U]
```
MapIterator uses parallelism goroutines to call f once for each element yielded
by iter. The returned iterator returns these results in the same order that iter
yielded them in.

This iterator, in contrast with most, must be consumed completely or it will
leak the goroutines.

If parallelism <= 0, uses GOMAXPROCS instead.

bufferSize is the size of the work buffer for each goroutine. A larger buffer
uses more memory but gives better throughput in the face of larger variance in
the processing time for f.

#### func  MapStream

```go
func MapStream[T any, U any](
	s stream.Stream[T],
	parallelism int,
	bufferSize int,
	f func(context.Context, T) (U, error),
) stream.Stream[U]
```
MapStream uses parallelism goroutines to call f once for each element yielded by
s. The returned stream returns these results in the same order that s yielded
them in.

If any call to f returns an error the context passed to invocations of f is
cancelled, no further calls to f are made, and the returned stream's Next
returns the first error encountered.

If parallelism <= 0, uses GOMAXPROCS instead.

bufferSize is the size of the work buffer for each goroutine. A larger buffer
uses more memory but gives better throughput in the face of larger variance in
the processing time for f.
