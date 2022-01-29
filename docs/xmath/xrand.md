# `package xrand`

```
import "github.com/bradenaw/juniper/xmath/xrand"
```

# Overview



# Index

<pre><a href="#Sample">func Sample(r *rand.Rand, n int, k int) []int</a></pre>
<pre><a href="#SampleIterator">func SampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T</a></pre>
<pre><a href="#SampleSlice">func SampleSlice[T any](r *rand.Rand, a []T, k int) []T</a></pre>
<pre><a href="#SampleStream">func SampleStream[T any](ctx context.Context, r *rand.Rand, s stream.Stream[T], k int) ([]T, error)</a></pre>
<pre><a href="#Shuffle">func Shuffle[T any](a []T)</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

## <a id="Sample"></a><pre>func <a href="#Sample">Sample</a>(r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, n int, k int) []int</pre>

Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).

If n < k, returns all ints in [0, n).

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


### Example 
```go
{
	r := rand.New(rand.NewSource(0))

	sample := xrand.Sample(r, 100, 5)

	fmt.Println(sample)

}
```

Output:
```text

```

## <a id="SampleIterator"></a><pre>func <a href="#SampleIterator">SampleIterator</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, iter <a href="../iterator.md#Iterator">iterator.Iterator</a>[T], k int) []T</pre>

SampleIterator pseudo-randomly picks k items uniformly without replacement from iter.

If iter yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


## <a id="SampleSlice"></a><pre>func <a href="#SampleSlice">SampleSlice</a>[T any](r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, a []T, k int) []T</pre>

SampleSlice pseudo-randomly picks k items uniformly without replacement from a.

If len(a) < k, returns all items in a.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


## <a id="SampleStream"></a><pre>func <a href="#SampleStream">SampleStream</a>[T any](ctx <a href="https://pkg.go.dev/context#Context">context.Context</a>, r *<a href="https://pkg.go.dev/math/rand#Rand">rand.Rand</a>, s <a href="../stream.md#Stream">stream.Stream</a>[T], k int) ([]T, error)</pre>

SampleStream pseudo-randomly picks k items uniformly without replacement from s.

If s yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired, the output should
be passed to Shuffle.


## <a id="Shuffle"></a><pre>func <a href="#Shuffle">Shuffle</a>[T any](a []T)</pre>

Shuffle pseudo-randomizes the order of a.


# Types

