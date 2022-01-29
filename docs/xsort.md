# `package xsort`

```
import "github.com/bradenaw/juniper/xsort"
```

# Overview

Package xsort contains extensions to the standard library package sort.


# Index

<pre><a href="#Equal">func Equal[T any](less Less[T], a T, b T) bool</a></pre>
<pre><a href="#Greater">func Greater[T any](less Less[T], a T, b T) bool</a></pre>
<pre><a href="#GreaterOrEqual">func GreaterOrEqual[T any](less Less[T], a T, b T) bool</a></pre>
<pre><a href="#LessOrEqual">func LessOrEqual[T any](less Less[T], a T, b T) bool</a></pre>
<pre><a href="#Merge">func Merge[T any](less Less[T], in ...iterator.Iterator[T]) iterator.Iterator[T]</a></pre>
<pre><a href="#MergeSlices">func MergeSlices[T any](less Less[T], out []T, in ...[]T) []T</a></pre>
<pre><a href="#MinK">func MinK[T any](less Less[T], iter iterator.Iterator[T], k int) []T</a></pre>
<pre><a href="#OrderedLess">func OrderedLess[T constraints.Ordered](a, b T) bool</a></pre>
<pre><a href="#Search">func Search[T any](x []T, less Less[T], item T) int</a></pre>
<pre><a href="#Slice">func Slice[T any](x []T, less Less[T])</a></pre>
<pre><a href="#SliceIsSorted">func SliceIsSorted[T any](x []T, less Less[T]) bool</a></pre>
<pre><a href="#SliceStable">func SliceStable[T any](x []T, less Less[T])</a></pre>
<pre><a href="#Less">type Less</a></pre>
<pre>    <a href="#Reverse">func Reverse[T any](less Less[T]) Less[T]</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

## <a id="Equal"></a><pre>func <a href="#Equal">Equal</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</pre>

Equal returns true if a == b according to less.


## <a id="Greater"></a><pre>func <a href="#Greater">Greater</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</pre>

Greater returns true if a > b according to less.


## <a id="GreaterOrEqual"></a><pre>func <a href="#GreaterOrEqual">GreaterOrEqual</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</pre>

LessOrEqual returns true if a >= b according to less.


## <a id="LessOrEqual"></a><pre>func <a href="#LessOrEqual">LessOrEqual</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</pre>

LessOrEqual returns true if a <= b according to less.


## <a id="Merge"></a><pre>func <a href="#Merge">Merge</a>[T any](less <a href="#Less">Less</a>[T], in ...) <a href="./iterator.md#Iterator">iterator.Iterator</a>[T]</pre>

Merge returns an iterator that yields all items from in in sorted order.

The results are undefined if the in iterators do not yield items in sorted order according to
less.

The time complexity of Next() is O(log(k)) where k is len(in).


### Example 
```go
{
	listOne := []string{"a", "f", "p", "x"}
	listTwo := []string{"b", "e", "o", "v"}
	listThree := []string{"s", "z"}

	merged := xsort.Merge(
		xsort.OrderedLess[string],
		iterator.Slice(listOne),
		iterator.Slice(listTwo),
		iterator.Slice(listThree),
	)

	fmt.Println(iterator.Collect(merged))

}
```

Output:
```text
[a b e f o p s v x z]

```

## <a id="MergeSlices"></a><pre>func <a href="#MergeSlices">MergeSlices</a>[T any](less <a href="#Less">Less</a>[T], out []T, in ...) []T</pre>

Merge merges the already-sorted slices of in. Optionally, a pre-allocated out slice can be
provided to store the result into.

The results are undefined if the in slices are not already sorted.

The time complexity is O(n * log(k)) where n is the total number of items and k is len(in).


### Example 
```go
{
	listOne := []string{"a", "f", "p", "x"}
	listTwo := []string{"b", "e", "o", "v"}
	listThree := []string{"s", "z"}

	merged := xsort.MergeSlices(
		xsort.OrderedLess[string],
		nil,
		listOne,
		listTwo,
		listThree,
	)

	fmt.Println(merged)

}
```

Output:
```text
[a b e f o p s v x z]

```

## <a id="MinK"></a><pre>func <a href="#MinK">MinK</a>[T any](less <a href="#Less">Less</a>[T], iter <a href="./iterator.md#Iterator">iterator.Iterator</a>[T], k int) []T</pre>

MinK returns the k minimum items according to less from iter in sorted order. If iter yields
fewer than k items, MinK returns all of them.


### Example 
```go
{
	a := []int{7, 4, 3, 8, 2, 1, 6, 9, 0, 5}

	iter := iterator.Slice(a)
	min3 := xsort.MinK(xsort.OrderedLess[int], iter, 3)
	fmt.Println(min3)

	iter = iterator.Slice(a)
	max3 := xsort.MinK(xsort.Reverse(xsort.OrderedLess[int]), iter, 3)
	fmt.Println(max3)

}
```

Output:
```text
[0 1 2]
[9 8 7]

```

## <a id="OrderedLess"></a><pre>func <a href="#OrderedLess">OrderedLess</a>[T <a href="https://pkg.go.dev/constraints#Ordered">constraints.Ordered</a>](a, b T) bool</pre>

OrderedLess is an implementation of Less for constraints.Ordered types by using the < operator.


## <a id="Search"></a><pre>func <a href="#Search">Search</a>[T any](x []T, less <a href="#Less">Less</a>[T], item T) int</pre>

Search searches for item in x, assumed sorted according to less, and returns the index. The
return value is the index to insert item at if it is not present (it could be len(a)).


### Example 
```go
{
	x := []string{"a", "f", "h", "i", "p", "z"}

	fmt.Println(xsort.Search(x, xsort.OrderedLess[string], "h"))
	fmt.Println(xsort.Search(x, xsort.OrderedLess[string], "k"))

}
```

Output:
```text
2
4

```

## <a id="Slice"></a><pre>func <a href="#Slice">Slice</a>[T any](x []T, less <a href="#Less">Less</a>[T])</pre>

Slice sorts x in-place using the given less function to compare items.

Follows the same rules as sort.Slice.


## <a id="SliceIsSorted"></a><pre>func <a href="#SliceIsSorted">SliceIsSorted</a>[T any](x []T, less <a href="#Less">Less</a>[T]) bool</pre>

SliceIsSorted returns true if x is in sorted order according to the given less function.

Follows the same rules as sort.SliceIsSorted.


## <a id="SliceStable"></a><pre>func <a href="#SliceStable">SliceStable</a>[T any](x []T, less <a href="#Less">Less</a>[T])</pre>

SliceStable stably sorts x in-place using the given less function to compare items.

Follows the same rules as sort.SliceStable.


# Types

## <a id="Less"></a><pre>type Less</pre>
```go
type Less[T any] func(a, b T) bool
```

Returns true if a is less than b. Must follow the same rules as sort.Interface.Less.


## <a id="Reverse"></a><pre>func Reverse[T any](less <a href="#Less">Less</a>[T]) <a href="#Less">Less</a>[T]</pre>

Reverse returns a Less that orders elements in the opposite order of the provided less.


