# `package xrand`

```
import "github.com/bradenaw/juniper/xmath/xrand"
```

## Overview



## Index

<samp><a href="#RSample">func RSample(r *rand.Rand, n int, k int) []int</a></samp>

<samp><a href="#RSampleIterator">func RSampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T</a></samp>

<samp><a href="#RSampleSlice">func RSampleSlice[T any](r *rand.Rand, a []T, k int) []T</a></samp>

<samp><a href="#RSampleStream">func RSampleStream[T any](ctx context.Context, r *rand.Rand, s stream.Stream[T], k int) ([]T, error)</a></samp>

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

<h3><a id="RSample"></a><samp>func <a href="#RSample">RSample</a>(r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, n int, k int) []int</samp></h3>

RSample pseudo-randomly picks k ints uniformly without replacement from [0, n).

If n < k, returns all ints in [0, n).

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="RSampleIterator"></a><samp>func <a href="#RSampleIterator">RSampleIterator</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, iter <a href="../iterator.html#Iterator">iterator.Iterator</a>[T], k int) []T</samp></h3>

RSampleIterator pseudo-randomly picks k items uniformly without replacement from iter.

If iter yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="RSampleSlice"></a><samp>func <a href="#RSampleSlice">RSampleSlice</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, a []T, k int) []T</samp></h3>

RSampleSlice pseudo-randomly picks k items uniformly without replacement from a.

If len(a) < k, returns all items in a.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="RSampleStream"></a><samp>func <a href="#RSampleStream">RSampleStream</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, s <a href="../stream.html#Stream">stream.Stream</a>[T], k int) ([]T, error)</samp></h3>

RSampleStream pseudo-randomly picks k items uniformly without replacement from s.

If s yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="RShuffle"></a><samp>func <a href="#RShuffle">RShuffle</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, a []T)</samp></h3>

RShuffle pseudo-randomizes the order of a.


<h3><a id="Sample"></a><samp>func <a href="#Sample">Sample</a>(n int, k int) []int</samp></h3>

Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).

If n < k, returns all ints in [0, n).

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


#### Example 
```go
{
	r := rand.New(rand.NewSource(0))

	sample := xrand.RSample(r, 100, 5)

	fmt.Println(sample)

}
```

Output:
```text
[45 71 88 93 60]
```
<h3><a id="SampleIterator"></a><samp>func <a href="#SampleIterator">SampleIterator</a>[T any](iter <a href="../iterator.html#Iterator">iterator.Iterator</a>[T], k int) []T</samp></h3>

SampleIterator pseudo-randomly picks k items uniformly without replacement from iter.

If iter yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="SampleSlice"></a><samp>func <a href="#SampleSlice">SampleSlice</a>[T any](a []T, k int) []T</samp></h3>

SampleSlice pseudo-randomly picks k items uniformly without replacement from a.

If len(a) < k, returns all items in a.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="SampleStream"></a><samp>func <a href="#SampleStream">SampleStream</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, s <a href="../stream.html#Stream">stream.Stream</a>[T], k int) ([]T, error)</samp></h3>

SampleStream pseudo-randomly picks k items uniformly without replacement from s.

If s yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


<h3><a id="Shuffle"></a><samp>func <a href="#Shuffle">Shuffle</a>[T any](a []T)</samp></h3>

Shuffle pseudo-randomizes the order of a.


## Types

This section is empty.

