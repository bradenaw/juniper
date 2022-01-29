# iterator
--
    import "."

package iterator allows iterating over sequences of values, for example the
contents of a container.

## Usage

#### func  Collect

```go
func Collect[T any](iter Iterator[T]) []T
```
Collect advances iter to the end and returns all of the items seen as a slice.

#### func  Equal

```go
func Equal[T comparable](iters ...Iterator[T]) bool
```
Equal returns true if the given iterators yield the same items in the same
order. Consumes the iterators.

#### func  Last

```go
func Last[T any](iter Iterator[T], n int) []T
```
Last consumes iter and returns the last n items. If iter yields fewer than n
items, Last returns all of them.

#### func  Reduce

```go
func Reduce[T any, U any](iter Iterator[T], initial U, f func(U, T) U) U
```
Reduce reduces iter to a single value using the reduction function f.

#### type Iterator

```go
type Iterator[T any] interface {
	// Next advances the iterator and returns the next item. Once the iterator is finished, the
	// first return is meaningless and the second return is false. Note that the final value of the
	// iterator has true in the second return, and it's the following call that returns false in the
	// second return.
	Next() (T, bool)
}
```

Iterator is used to iterate over a sequence of values.

Iterators are lazy, meaning they do no work until a call to Next().

Iterators do not need to be fully consumed, callers may safely abandon an
iterator before Next returns false.

#### func  Chan

```go
func Chan[T any](c <-chan T) Iterator[T]
```
Chan returns an Iterator that yields the values received on c.

#### func  Chunk

```go
func Chunk[T any](iter Iterator[T], chunkSize int) Iterator[[]T]
```
Chunk returns an iterator over non-overlapping chunks of size chunkSize. The
last chunk will be smaller than chunkSize if the iterator does not contain an
even multiple.

#### func  Compact

```go
func Compact[T comparable](iter Iterator[T]) Iterator[T]
```
Compact elides adjacent duplicates from iter.

#### func  CompactFunc

```go
func CompactFunc[T any](iter Iterator[T], eq func(T, T) bool) Iterator[T]
```
CompactFunc elides adjacent duplicates from iter, using eq to determine
duplicates.

#### func  Counter

```go
func Counter(n int) Iterator[int]
```
Counter returns an iterator that counts up from 0, yielding n items.

The following are equivalent:

    for i := 0; i < n; i++ {
      fmt.Println(n)
    }

    iter := iterator.Counter(n)
    for {
      item, ok := iter.Next()
      if !ok {
        break
      }
      fmt.Println(item)
    }

#### func  Filter

```go
func Filter[T any](iter Iterator[T], keep func(T) bool) Iterator[T]
```
Filter returns an iterator that yields only the items from iter for which keep
returns true.

#### func  First

```go
func First[T any](iter Iterator[T], n int) Iterator[T]
```
First returns an iterator that yields the first n items from iter.

#### func  Flatten

```go
func Flatten[T any](iter Iterator[Iterator[T]]) Iterator[T]
```
Flatten returns an iterator that yields all items from all iterators yielded by
iter.

#### func  Join

```go
func Join[T any](iters ...Iterator[T]) Iterator[T]
```
Join returns an Iterator that returns all elements of iters[0], then all
elements of iters[1], and so on.

#### func  Map

```go
func Map[T any, U any](iter Iterator[T], f func(t T) U) Iterator[U]
```
Map transforms the results of iter using the conversion f.

#### func  Repeat

```go
func Repeat[T any](item T, n int) Iterator[T]
```
Repeat returns an iterator that yields item n times.

#### func  Runs

```go
func Runs[T any](iter Iterator[T], same func(a, b T) bool) Iterator[Iterator[T]]
```
Runs returns an iterator of iterators. The inner iterators yield contiguous
elements from iter such that same(a, b) returns true for any a and b in the run.

The inner iterator should be drained before calling Next on the outer iterator.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then
same(a, c) must also.

#### func  Slice

```go
func Slice[T any](s []T) Iterator[T]
```
Slice returns an iterator over the elements of s.

#### func  While

```go
func While[T any](iter Iterator[T], f func(T) bool) Iterator[T]
```
While returns an iterator that terminates before the first item from iter for
which f returns false.

#### type Peekable

```go
type Peekable[T any] interface {
	Iterator[T]
	// Peek returns the next item of the iterator if there is one without consuming it.
	//
	// If Peek returns a value, the next call to Next will return the same value.
	Peek() (T, bool)
}
```

Peekable allows viewing the next item from an iterator without consuming it.

#### func  WithPeek

```go
func WithPeek[T any](iter Iterator[T]) Peekable[T]
```
WithPeek returns iter with a Peek() method attached.
