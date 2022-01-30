# `package parallel`

```
import "github.com/bradenaw/juniper/parallel"
```

## Overview

Package parallel provides primitives for running tasks in parallel.


## Index

<samp><a href="#Do">func Do(
	parallelism int,
	n int,
	f func(i int),
)</a></samp>

<samp><a href="#DoContext">func DoContext(
	ctx context.Context,
	parallelism int,
	n int,
	f func(ctx context.Context, i int) error,
) error</a></samp>

<samp><a href="#Map">func Map[T any, U any](
	parallelism int,
	in []T,
	f func(in T) U,
) []U</a></samp>

<samp><a href="#MapContext">func MapContext[T any, U any](
	ctx context.Context,
	parallelism int,
	in []T,
	f func(ctx context.Context, in T) (U, error),
) ([]U, error)</a></samp>

<samp><a href="#MapIterator">func MapIterator[T any, U any](
	iter iterator.Iterator[T],
	parallelism int,
	bufferSize int,
	f func(T) U,
) iterator.Iterator[U]</a></samp>

<samp><a href="#MapStream">func MapStream[T any, U any](
	s stream.Stream[T],
	parallelism int,
	bufferSize int,
	f func(context.Context, T) (U, error),
) stream.Stream[U]</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="Do"></a><samp>func <a href="#Do">Do</a>(parallelism int, n int, f func(i int))</samp></h3>

Do calls f from parallelism goroutines n times, providing each invocation a unique i in [0, n).

If parallelism <= 0, uses GOMAXPROCS instead.


<h3><a id="DoContext"></a><samp>func <a href="#DoContext">DoContext</a>(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, parallelism int, n int, f func(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, i int) error) error</samp></h3>

DoContext calls f from parallelism goroutines n times, providing each invocation a unique i in
[0, n).

If any call to f returns an error the context passed to invocations of f is cancelled, no further
calls to f are made, and Do returns the first error encountered.

If parallelism <= 0, uses GOMAXPROCS instead.


<h3><a id="Map"></a><samp>func <a href="#Map">Map</a>[T any, U any](parallelism int, in []T, f func(in T) U) []U</samp></h3>

Map uses parallelism goroutines to call f once for each element of in. out[i] is the
result of f for in[i].

If parallelism <= 0, uses GOMAXPROCS instead.


<h3><a id="MapContext"></a><samp>func <a href="#MapContext">MapContext</a>[T any, U any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, parallelism int, in []T, f func(ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, in T) (U, error)) ([]U, error)</samp></h3>

MapContext uses parallelism goroutines to call f once for each element of in. out[i] is the
result of f for in[i].

If any call to f returns an error the context passed to invocations of f is cancelled, no further
calls to f are made, and Map returns the first error encountered.

If parallelism <= 0, uses GOMAXPROCS instead.


<h3><a id="MapIterator"></a><samp>func <a href="#MapIterator">MapIterator</a>[T any, U any](iter <a href="./iterator.html#Iterator">iterator.Iterator</a>[T], parallelism int, bufferSize int, f func(T) U) <a href="./iterator.html#Iterator">iterator.Iterator</a>[U]</samp></h3>

MapIterator uses parallelism goroutines to call f once for each element yielded by iter. The
returned iterator returns these results in the same order that iter yielded them in.

This iterator, in contrast with most, must be consumed completely or it will leak the goroutines.

If parallelism <= 0, uses GOMAXPROCS instead.

bufferSize is the size of the work buffer for each goroutine. A larger buffer uses more memory
but gives better throughput in the face of larger variance in the processing time for f.


<h3><a id="MapStream"></a><samp>func <a href="#MapStream">MapStream</a>[T any, U any](s <a href="./stream.html#Stream">stream.Stream</a>[T], parallelism int, bufferSize int, f func(<a href="https://pkg.go.dev/context#Context">context.Context</a>, T) (U, error)) <a href="./stream.html#Stream">stream.Stream</a>[U]</samp></h3>

MapStream uses parallelism goroutines to call f once for each element yielded by s. The returned
stream returns these results in the same order that s yielded them in.

If any call to f returns an error the context passed to invocations of f is cancelled, no further
calls to f are made, and the returned stream's Next returns the first error encountered.

If parallelism <= 0, uses GOMAXPROCS instead.

bufferSize is the size of the work buffer for each goroutine. A larger buffer uses more memory
but gives better throughput in the face of larger variance in the processing time for f.


## Types

