# slices
--
    import "."


## Usage

#### func  All

```go
func All[T any](x []T, f func(T) bool) bool
```
All returns true if f(x[i]) returns true for all i. Trivially, returns true if x
is empty.

#### func  Any

```go
func Any[T any](x []T, f func(T) bool) bool
```
Any returns true if f(x[i]) returns true for any i. Trivially, returns false if
x is empty.

#### func  Chunk

```go
func Chunk[T any](x []T, chunkSize int) [][]T
```
Chunk returns non-overlapping chunks of x. The last chunk will be smaller than
chunkSize if len(x) is not a multiple of chunkSize.

Returns an empty slice if len(x)==0. Panics if chunkSize <= 0.

#### func  Clear

```go
func Clear[T any](x []T)
```
Clear fills x with the zero value of T.

#### func  Clone

```go
func Clone[T any](x []T) []T
```
Clone creates a new slice and copies the elements of x into it.

#### func  Compact

```go
func Compact[T comparable](x []T) []T
```
Compact removes adjacent duplicates from x in-place and returns the modified
slice.

#### func  CompactFunc

```go
func CompactFunc[T any](x []T, eq func(T, T) bool) []T
```
CompactFunc removes adjacent duplicates from x in-place, preserving the first
occurrence, using the supplied eq function and returns the modified slice.

#### func  Count

```go
func Count[T comparable](a []T, item T) int
```
Count returns the number of times item appears in a.

#### func  CountFunc

```go
func CountFunc[T any](a []T, f func(T) bool) int
```
Count returns the number of items in a for which f returns true.

#### func  Equal

```go
func Equal[T comparable](a, b []T) bool
```
Equal returns true if a and b contain the same items in the same order.

#### func  Fill

```go
func Fill[T any](a []T, x T)
```
Fill fills a with copies of x.

#### func  Filter

```go
func Filter[T any](x []T, keep func(t T) bool) []T
```
Filter filters the contents of x to only those for which keep() returns true.
This is done in-place and so modifies the contents of x. The modified slice is
returned.

#### func  Flatten

```go
func Flatten[T any](x [][]T) []T
```
Flatten returns a slice containing all of the elements of all elements of x.

#### func  Grow

```go
func Grow[T any](x []T, n int) []T
```
Grow grows x's capacity by reallocating, if necessary, to fit n more elements
and returns the modified slice. This does not change the length of x. After
Grow(x, n), the following n append()s to x will not need to reallocate.

#### func  Index

```go
func Index[T comparable](a []T, item T) int
```
Index returns the first index of item in a, or -1 if item is not in a.

#### func  IndexFunc

```go
func IndexFunc[T any](a []T, f func(T) bool) int
```
Index returns the first index in a for which f(a[i]) returns true, or -1 if
there are no such items.

#### func  Insert

```go
func Insert[T any](x []T, idx int, values ...T) []T
```
Insert inserts the given values starting at index idx, shifting elements after
idx to the right and growing the slice to make room. Insert will expand the
length of the slice up to its capacity if it can, if this isn't desired then x
should be resliced to have capacity equal to its length:

    x[:len(x):len(x)]

The time cost is O(n+m) where n is len(values) and m is len(x[idx:]).

#### func  Join

```go
func Join[T any](in ...[]T) []T
```
Join joins together the contents of each in.

#### func  LastIndex

```go
func LastIndex[T comparable](a []T, item T) int
```
LastIndex returns the last index of item in a, or -1 if item is not in a.

#### func  LastIndexFunc

```go
func LastIndexFunc[T any](a []T, f func(T) bool) int
```
LastIndexFunc returns the last index in a for which f(a[i]) returns true, or -1
if there are no such items.

#### func  Map

```go
func Map[T any, U any](x []T, f func(T) U) []U
```
Map creates a new slice by applying f to each element of x.

#### func  Partition

```go
func Partition[T any](x []T, f func(t T) bool)
```
Partition moves elements of x such that all elements for which f returns false
are at the beginning and all elements for which f returns true are at the end.
It makes no other guarantees about the final order of elements.

#### func  Reduce

```go
func Reduce[T any, U any](x []T, initial U, f func(U, T) U) U
```
Reduce reduces x to a single value using the reduction function f.

#### func  Remove

```go
func Remove[T any](x []T, idx int, n int) []T
```
Remove removes n elements from x starting at index idx and returns the modified
slice. This requires shifting the elements after the removed elements over, and
so its cost is linear in the number of elements shifted.

#### func  Repeat

```go
func Repeat[T any](x T, n int) []T
```
Repeat returns a slice with length n where every item is x.

#### func  Reverse

```go
func Reverse[T any](x []T)
```
Reverse reverses the elements of x in place.

#### func  Runs

```go
func Runs[T any](x []T, same func(a, b T) bool) [][]T
```
Runs returns a slice of slices. The inner slices are contiguous runs of elements
from x such that same(a, b) returns true for any a and b in the run.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then
same(a, c) must also.

The returned slices use the same underlying array as x.

#### func  Shrink

```go
func Shrink[T any](x []T, n int) []T
```
Shrink shrinks x's capacity by reallocating, if necessary, so that cap(x) <=
len(x) + n.

#### func  Unique

```go
func Unique[T comparable](x []T) []T
```
Unique removes duplicates from x in-place, preserving order, and returns the
modified slice.

Compact is more efficient if duplicates are already adjacent in x, for example
if x is in sorted order.
