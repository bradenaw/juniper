# xrand
--
    import "."


## Usage

#### func  Sample

```go
func Sample(r *rand.Rand, n int, k int) []int
```
Sample pseudo-randomly picks k ints uniformly without replacement from [0, n).

If n < k, returns all ints in [0, n).

The output is not in any particular order. If a pseudo-random order is desired,
the output should be passed to Shuffle.

#### func  SampleIterator

```go
func SampleIterator[T any](r *rand.Rand, iter iterator.Iterator[T], k int) []T
```
SampleIterator pseudo-randomly picks k items uniformly without replacement from
iter.

If iter yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired,
the output should be passed to Shuffle.

#### func  SampleSlice

```go
func SampleSlice[T any](r *rand.Rand, a []T, k int) []T
```
SampleSlice pseudo-randomly picks k items uniformly without replacement from a.

If len(a) < k, returns all items in a.

The output is not in any particular order. If a pseudo-random order is desired,
the output should be passed to Shuffle.

#### func  SampleStream

```go
func SampleStream[T any](ctx context.Context, r *rand.Rand, s stream.Stream[T], k int) ([]T, error)
```
SampleStream pseudo-randomly picks k items uniformly without replacement from s.

If s yields fewer than k items, returns all of them.

The output is not in any particular order. If a pseudo-random order is desired,
the output should be passed to Shuffle.

#### func  Shuffle

```go
func Shuffle[T any](a []T)
```
Shuffle pseudo-randomizes the order of a.
