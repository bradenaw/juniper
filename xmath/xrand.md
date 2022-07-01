# `package xrand`

```
import "github.com/bradenaw/juniper/xmath/xrand"
```

## Overview

Package xrand contains extensions to the standard library package math/rand.


## Index

<samp><a href="#RSample">func RSample(r *rand.Rand, n int, k int) []int</a></samp>

<samp><a href="#RSampleIterator">func RSampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T</a></samp>

<samp><a href="#RSampleSlice">func RSampleSlice[T any](r *rand.Rand, a []T, k int) []T</a></samp>

<samp><a href="#RSampleStream">func RSampleStream[T any](
	ctx context.Context,
	r *rand.Rand,
	s stream.Stream[T],
	k int,
) ([]T, error)</a></samp>

<samp><a href="#RShuffle">func RShuffle[T any](r *rand.Rand, a []T)</a></samp>

<samp><a href="#Sample">func Sample(n int, k int) []int</a></samp>

<samp><a href="#SampleIterator">func SampleIterator[T any](iter iterator.Iterator[T], k int) []T</a></samp>

<samp><a href="#SampleSlice">func SampleSlice[T any](a []T, k int) []T</a></samp>

<samp><a href="#SampleStream">func SampleStream[T any](ctx context.Context, s stream.Stream[T], k int) ([]T, error)</a></samp>

<samp><a href="#Shuffle">func Shuffle[T any](a []T)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="RSample"></a><samp>func <a href="#RSample">RSample</a>(r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, n int, k int) []int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L55">src</a></small></sub></h3>

RSample pseudo-randomly picks k ints uniformly without replacement from [0, n).

If n < k, returns all ints in [0, n).

Requires O(k) time and space.


<h3><a id="RSampleIterator"></a><samp>func <a href="#RSampleIterator">RSampleIterator</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, iter <a href="../iterator.html#Iterator">iterator.Iterator</a>[T], k int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L92">src</a></small></sub></h3>

RSampleIterator pseudo-randomly picks k items uniformly without replacement from iter.

If iter yields fewer than k items, returns all of them.

Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
linear in the length of iter but only O(k) extra space.


<h3><a id="RSampleSlice"></a><samp>func <a href="#RSampleSlice">RSampleSlice</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, a []T, k int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L200">src</a></small></sub></h3>

RSampleSlice pseudo-randomly picks k items uniformly without replacement from a.

If len(a) < k, returns all items in a.

Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
and space.


<h3><a id="RSampleStream"></a><samp>func <a href="#RSampleStream">RSampleStream</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, s <a href="../stream.html#Stream">stream.Stream</a>[T], k int) ([]T, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L139">src</a></small></sub></h3>

RSampleStream pseudo-randomly picks k items uniformly without replacement from s.

If s yields fewer than k items, returns all of them.

Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
linear in the length of s but only O(k) extra space.


<h3><a id="RShuffle"></a><samp>func <a href="#RShuffle">RShuffle</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, a []T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L31">src</a></small></sub></h3>

RShuffle pseudo-randomizes the order of a.


<h3><a id="Sample"></a><samp>func <a href="#Sample">Sample</a>(n int, k int) []int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L46">src</a></small></sub></h3>

Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).

If n < k, returns all ints in [0, n).

Requires O(k) time and space.


#### Example 
```go
{
	r := rand.New(rand.NewSource(0))

	sample := xrand.RSample(r, 100, 5)

	for _, x := range sample {
		fmt.Println(x)
	}

}
```

Unordered output:
```text
45
71
88
93
60
```
<h3><a id="SampleIterator"></a><samp>func <a href="#SampleIterator">SampleIterator</a>[T any](iter <a href="../iterator.html#Iterator">iterator.Iterator</a>[T], k int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L82">src</a></small></sub></h3>

SampleIterator pseudo-randomly picks k items uniformly without replacement from iter.

If iter yields fewer than k items, returns all of them.

Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
linear in the length of iter but only O(k) extra space.


<h3><a id="SampleSlice"></a><samp>func <a href="#SampleSlice">SampleSlice</a>[T any](a []T, k int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L190">src</a></small></sub></h3>

SampleSlice pseudo-randomly picks k items uniformly without replacement from a.

If len(a) < k, returns all items in a.

Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses O(k) time
and space.


<h3><a id="SampleStream"></a><samp>func <a href="#SampleStream">SampleStream</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="../stream.html#Stream">stream.Stream</a>[T], k int) ([]T, error)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L129">src</a></small></sub></h3>

SampleStream pseudo-randomly picks k items uniformly without replacement from s.

If s yields fewer than k items, returns all of them.

Uses a reservoir sample (https://en.wikipedia.org/wiki/Reservoir_sampling), which uses time
linear in the length of s but only O(k) extra space.


<h3><a id="Shuffle"></a><samp>func <a href="#Shuffle">Shuffle</a>[T any](a []T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xmath/xrand/xrand.go#L26">src</a></small></sub></h3>

Shuffle pseudo-randomizes the order of a.


## Types

This section is empty.

