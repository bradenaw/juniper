# xsort
--
    import "."

Package xsort contains extensions to the standard library package sort.

## Usage

#### func  Equal

```go
func Equal[T any](less Less[T], a T, b T) bool
```
Equal returns true if a == b according to less.

#### func  Greater

```go
func Greater[T any](less Less[T], a T, b T) bool
```
Greater returns true if a > b according to less.

#### func  GreaterOrEqual

```go
func GreaterOrEqual[T any](less Less[T], a T, b T) bool
```
LessOrEqual returns true if a >= b according to less.

#### func  LessOrEqual

```go
func LessOrEqual[T any](less Less[T], a T, b T) bool
```
LessOrEqual returns true if a <= b according to less.

#### func  Merge

```go
func Merge[T any](less Less[T], in ...iterator.Iterator[T]) iterator.Iterator[T]
```
Merge returns an iterator that yields all items from in in sorted order.

The results are undefined if the in iterators do not yield items in sorted order
according to less.

The time complexity of Next() is O(log(k)) where k is len(in).

#### func  MergeSlices

```go
func MergeSlices[T any](less Less[T], out []T, in ...[]T) []T
```
Merge merges the already-sorted slices of in. Optionally, a pre-allocated out
slice can be provided to store the result into.

The results are undefined if the in slices are not already sorted.

The time complexity is O(n * log(k)) where n is the total number of items and k
is len(in).

#### func  MinK

```go
func MinK[T any](less Less[T], iter iterator.Iterator[T], k int) []T
```
MinK returns the k minimum items according to less from iter in sorted order. If
iter yields fewer than k items, MinK returns all of them.

#### func  OrderedLess

```go
func OrderedLess[T constraints.Ordered](a, b T) bool
```
OrderedLess is an implementation of Less for constraints.Ordered types by using
the < operator.

#### func  Search

```go
func Search[T any](x []T, less Less[T], item T) int
```
Search searches for item in x, assumed sorted according to less, and returns the
index. The return value is the index to insert item at if it is not present (it
could be len(a)).

#### func  Slice

```go
func Slice[T any](x []T, less Less[T])
```
Slice sorts x in-place using the given less function to compare items.

Follows the same rules as sort.Slice.

#### func  SliceIsSorted

```go
func SliceIsSorted[T any](x []T, less Less[T]) bool
```
SliceIsSorted returns true if x is in sorted order according to the given less
function.

Follows the same rules as sort.SliceIsSorted.

#### func  SliceStable

```go
func SliceStable[T any](x []T, less Less[T])
```
SliceStable stably sorts x in-place using the given less function to compare
items.

Follows the same rules as sort.SliceStable.

#### type Less

```go
type Less[T any] func(a, b T) bool
```

Returns true if a is less than b. Must follow the same rules as
sort.Interface.Less.

#### func  Reverse

```go
func Reverse[T any](less Less[T]) Less[T]
```
Reverse returns a Less that orders elements in the opposite order of the
provided less.
