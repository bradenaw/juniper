# `package xsort`

```
import "github.com/bradenaw/juniper/xsort"
```

# Overview

Package xsort contains extensions to the standard library package sort.


# Index

<samp><a href="#Equal">func Equal[T any](less Less[T], a T, b T) bool</a></samp>
<samp><a href="#Greater">func Greater[T any](less Less[T], a T, b T) bool</a></samp>
<samp><a href="#GreaterOrEqual">func GreaterOrEqual[T any](less Less[T], a T, b T) bool</a></samp>
<samp><a href="#LessOrEqual">func LessOrEqual[T any](less Less[T], a T, b T) bool</a></samp>
<samp><a href="#Merge">func Merge[T any](less Less[T], in ...iterator.Iterator[T]) iterator.Iterator[T]</a></samp>
<samp><a href="#MergeSlices">func MergeSlices[T any](less Less[T], out []T, in ...[]T) []T</a></samp>
<samp><a href="#MinK">func MinK[T any](less Less[T], iter iterator.Iterator[T], k int) []T</a></samp>
<samp><a href="#OrderedLess">func OrderedLess[T constraints.Ordered](a, b T) bool</a></samp>
<samp><a href="#Search">func Search[T any](x []T, less Less[T], item T) int</a></samp>
<samp><a href="#Slice">func Slice[T any](x []T, less Less[T])</a></samp>
<samp><a href="#SliceIsSorted">func SliceIsSorted[T any](x []T, less Less[T]) bool</a></samp>
<samp><a href="#SliceStable">func SliceStable[T any](x []T, less Less[T])</a></samp>
<samp><a href="#Less">type Less</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Reverse">func Reverse[T any](less Less[T]) Less[T]</a></samp>


# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="Equal"></a><samp>func <a href="#Equal">Equal</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</samp></h2>

Equal returns true if a == b according to less.


<h2><a id="Greater"></a><samp>func <a href="#Greater">Greater</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</samp></h2>

Greater returns true if a > b according to less.


<h2><a id="GreaterOrEqual"></a><samp>func <a href="#GreaterOrEqual">GreaterOrEqual</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</samp></h2>

LessOrEqual returns true if a >= b according to less.


<h2><a id="LessOrEqual"></a><samp>func <a href="#LessOrEqual">LessOrEqual</a>[T any](less <a href="#Less">Less</a>[T], a T, b T) bool</samp></h2>

LessOrEqual returns true if a <= b according to less.


<h2><a id="Merge"></a><samp>func <a href="#Merge">Merge</a>[T any](less <a href="#Less">Less</a>[T], in ...) <a href="./iterator.md#Iterator">iterator.Iterator</a>[T]</samp></h2>

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
<h2><a id="MergeSlices"></a><samp>func <a href="#MergeSlices">MergeSlices</a>[T any](less <a href="#Less">Less</a>[T], out []T, in ...) []T</samp></h2>

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
<h2><a id="MinK"></a><samp>func <a href="#MinK">MinK</a>[T any](less <a href="#Less">Less</a>[T], iter <a href="./iterator.md#Iterator">iterator.Iterator</a>[T], k int) []T</samp></h2>

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
<h2><a id="OrderedLess"></a><samp>func <a href="#OrderedLess">OrderedLess</a>[T <a href="https://pkg.go.dev/constraints#Ordered">constraints.Ordered</a>](a, b T) bool</samp></h2>

OrderedLess is an implementation of Less for constraints.Ordered types by using the < operator.


<h2><a id="Search"></a><samp>func <a href="#Search">Search</a>[T any](x []T, less <a href="#Less">Less</a>[T], item T) int</samp></h2>

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
<h2><a id="Slice"></a><samp>func <a href="#Slice">Slice</a>[T any](x []T, less <a href="#Less">Less</a>[T])</samp></h2>

Slice sorts x in-place using the given less function to compare items.

Follows the same rules as sort.Slice.


<h2><a id="SliceIsSorted"></a><samp>func <a href="#SliceIsSorted">SliceIsSorted</a>[T any](x []T, less <a href="#Less">Less</a>[T]) bool</samp></h2>

SliceIsSorted returns true if x is in sorted order according to the given less function.

Follows the same rules as sort.SliceIsSorted.


<h2><a id="SliceStable"></a><samp>func <a href="#SliceStable">SliceStable</a>[T any](x []T, less <a href="#Less">Less</a>[T])</samp></h2>

SliceStable stably sorts x in-place using the given less function to compare items.

Follows the same rules as sort.SliceStable.


# Types

<h2><a id="Less"></a><samp>type Less</samp></h2>
```go
type Less[T any] func(a, b T) bool
```

Returns true if a is less than b. Must follow the same rules as sort.Interface.Less.


<h2><a id="Reverse"></a><samp>func Reverse[T any](less <a href="#Less">Less</a>[T]) <a href="#Less">Less</a>[T]</samp></h2>

Reverse returns a Less that orders elements in the opposite order of the provided less.


