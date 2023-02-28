# `package xslices`

```
import "github.com/bradenaw/juniper/xslices"
```

## Overview

Package xslices contains utilities for working with slices of arbitrary types.


## Index

<samp><a href="#All">func All[T any](s []T, f func(T) bool) bool</a></samp>

<samp><a href="#Any">func Any[T any](s []T, f func(T) bool) bool</a></samp>

<samp><a href="#Chunk">func Chunk[T any](s []T, chunkSize int) [][]T</a></samp>

<samp><a href="#Clear">func Clear[T any](s []T)</a></samp>

<samp><a href="#Clone">func Clone[T any](s []T) []T</a></samp>

<samp><a href="#Compact">func Compact[T comparable](s []T) []T</a></samp>

<samp><a href="#CompactFunc">func CompactFunc[T any](s []T, eq func(T, T) bool) []T</a></samp>

<samp><a href="#CompactInPlace">func CompactInPlace[T comparable](s []T) []T</a></samp>

<samp><a href="#CompactInPlaceFunc">func CompactInPlaceFunc[T any](s []T, eq func(T, T) bool) []T</a></samp>

<samp><a href="#Count">func Count[T comparable](s []T, x T) int</a></samp>

<samp><a href="#CountFunc">func CountFunc[T any](s []T, f func(T) bool) int</a></samp>

<samp><a href="#Equal">func Equal[T comparable](a, b []T) bool</a></samp>

<samp><a href="#EqualFunc">func EqualFunc[T any](a, b []T, eq func(T, T) bool) bool</a></samp>

<samp><a href="#Fill">func Fill[T any](s []T, x T)</a></samp>

<samp><a href="#Filter">func Filter[T any](s []T, keep func(t T) bool) []T</a></samp>

<samp><a href="#FilterInPlace">func FilterInPlace[T any](s []T, keep func(t T) bool) []T</a></samp>

<samp><a href="#Group">func Group[T any, U comparable](s []T, f func(T) U) map[U][]T</a></samp>

<samp><a href="#Grow">func Grow[T any](s []T, n int) []T</a></samp>

<samp><a href="#Index">func Index[T comparable](s []T, x T) int</a></samp>

<samp><a href="#IndexFunc">func IndexFunc[T any](s []T, f func(T) bool) int</a></samp>

<samp><a href="#Insert">func Insert[T any](s []T, idx int, values ...T) []T</a></samp>

<samp><a href="#Join">func Join[T any](in ...[]T) []T</a></samp>

<samp><a href="#LastIndex">func LastIndex[T comparable](s []T, x T) int</a></samp>

<samp><a href="#LastIndexFunc">func LastIndexFunc[T any](s []T, f func(T) bool) int</a></samp>

<samp><a href="#Map">func Map[T any, U any](s []T, f func(T) U) []U</a></samp>

<samp><a href="#Partition">func Partition[T any](s []T, f func(t T) bool) int</a></samp>

<samp><a href="#Reduce">func Reduce[T any, U any](s []T, initial U, f func(U, T) U) U</a></samp>

<samp><a href="#Remove">func Remove[T any](s []T, idx int, n int) []T</a></samp>

<samp><a href="#RemoveUnordered">func RemoveUnordered[T any](s []T, idx int, n int) []T</a></samp>

<samp><a href="#Repeat">func Repeat[T any](s T, n int) []T</a></samp>

<samp><a href="#Reverse">func Reverse[T any](s []T)</a></samp>

<samp><a href="#Runs">func Runs[T any](s []T, same func(a, b T) bool) [][]T</a></samp>

<samp><a href="#Shrink">func Shrink[T any](s []T, n int) []T</a></samp>

<samp><a href="#Unique">func Unique[T comparable](s []T) []T</a></samp>

<samp><a href="#UniqueInPlace">func UniqueInPlace[T comparable](s []T) []T</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="All"></a><samp>func <a href="#All">All</a>[T any](s []T, f func(T) bool) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L5">src</a></small></sub></h3>

All returns true if f(s[i]) returns true for all i. Trivially, returns true if s is empty.


#### Example 
```go
{
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	allOdd := xslices.All([]int{1, 3, 5}, isOdd)
	fmt.Println(allOdd)

	allOdd = xslices.All([]int{1, 3, 6}, isOdd)
	fmt.Println(allOdd)

}
```

Output:
```text
true
false
```
<h3><a id="Any"></a><samp>func <a href="#Any">Any</a>[T any](s []T, f func(T) bool) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L15">src</a></small></sub></h3>

Any returns true if f(s[i]) returns true for any i. Trivially, returns false if s is empty.


#### Example 
```go
{
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	anyOdd := xslices.Any([]int{2, 3, 4}, isOdd)
	fmt.Println(anyOdd)

	anyOdd = xslices.Any([]int{2, 4, 6}, isOdd)
	fmt.Println(anyOdd)

}
```

Output:
```text
true
false
```
<h3><a id="Chunk"></a><samp>func <a href="#Chunk">Chunk</a>[T any](s []T, chunkSize int) [][]T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L28">src</a></small></sub></h3>

Chunk returns non-overlapping chunks of s. The last chunk will be smaller than chunkSize if
len(s) is not a multiple of chunkSize.

Returns an empty slice if len(s)==0. Panics if chunkSize <= 0.


#### Example 
```go
{
	s := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	chunks := xslices.Chunk(s, 3)
	fmt.Println(chunks)

}
```

Output:
```text
[[a b c] [d e f] [g h]]
```
<h3><a id="Clear"></a><samp>func <a href="#Clear">Clear</a>[T any](s []T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L42">src</a></small></sub></h3>

Clear fills s with the zero value of T.


#### Example 
```go
{
	s := []int{1, 2, 3}
	xslices.Clear(s)
	fmt.Println(s)

}
```

Output:
```text
[0 0 0]
```
<h3><a id="Clone"></a><samp>func <a href="#Clone">Clone</a>[T any](s []T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L48">src</a></small></sub></h3>

Clone creates a new slice and copies the elements of s into it.


#### Example 
```go
{
	s := []int{1, 2, 3}
	cloned := xslices.Clone(s)
	fmt.Println(cloned)

}
```

Output:
```text
[1 2 3]
```
<h3><a id="Compact"></a><samp>func <a href="#Compact">Compact</a>[T comparable](s []T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L56">src</a></small></sub></h3>

Compact returns a slice containing only the first item from each contiguous run of the same item.

For example, this can be used to remove duplicates more cheaply than Unique when the slice is
already in sorted order.


#### Example 
```go
{
	s := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := xslices.Compact(s)
	fmt.Println(compacted)

}
```

Output:
```text
[a b c a]
```
<h3><a id="CompactFunc"></a><samp>func <a href="#CompactFunc">CompactFunc</a>[T any](s []T, eq func(T, T) bool) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L74">src</a></small></sub></h3>

CompactFunc returns a slice containing only the first item from each contiguous run of items for
which eq returns true.


#### Example 
```go
{
	s := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := xslices.CompactFunc(s, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

}
```

Output:
```text
[bank ghost yaw]
```
<h3><a id="CompactInPlace"></a><samp>func <a href="#CompactInPlace">CompactInPlace</a>[T comparable](s []T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L66">src</a></small></sub></h3>

CompactInPlace returns a slice containing only the first item from each contiguous run of the
same item. This is done in-place and so modifies the contents of s. The modified slice is
returned.

For example, this can be used to remove duplicates more cheaply than Unique when the slice is
already in sorted order.


#### Example 
```go
{
	s := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := xslices.CompactInPlace(s)
	fmt.Println(compacted)

}
```

Output:
```text
[a b c a]
```
<h3><a id="CompactInPlaceFunc"></a><samp>func <a href="#CompactInPlaceFunc">CompactInPlaceFunc</a>[T any](s []T, eq func(T, T) bool) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L81">src</a></small></sub></h3>

CompactInPlaceFunc returns a slice containing only the first item from each contiguous run of
items for which eq returns true. This is done in-place and so modifies the contents of s. The
modified slice is returned.


#### Example 
```go
{
	s := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := xslices.CompactInPlaceFunc(s, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

}
```

Output:
```text
[bank ghost yaw]
```
<h3><a id="Count"></a><samp>func <a href="#Count">Count</a>[T comparable](s []T, x T) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L97">src</a></small></sub></h3>

Count returns the number of times x appears in s.


#### Example 
```go
{
	s := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.Count(s, "a"))

}
```

Output:
```text
3
```
<h3><a id="CountFunc"></a><samp>func <a href="#CountFunc">CountFunc</a>[T any](s []T, f func(T) bool) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L102">src</a></small></sub></h3>

Count returns the number of items in s for which f returns true.


<h3><a id="Equal"></a><samp>func <a href="#Equal">Equal</a>[T comparable](a, b []T) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L113">src</a></small></sub></h3>

Equal returns true if a and b contain the same items in the same order.


#### Example 
```go
{
	x := []string{"a", "b", "c"}
	y := []string{"a", "b", "c"}
	z := []string{"a", "b", "d"}

	fmt.Println(xslices.Equal(x, y))
	fmt.Println(xslices.Equal(x[:2], y))
	fmt.Println(xslices.Equal(z, y))

}
```

Output:
```text
true
false
false
```
<h3><a id="EqualFunc"></a><samp>func <a href="#EqualFunc">EqualFunc</a>[T any](a, b []T, eq func(T, T) bool) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L126">src</a></small></sub></h3>

EqualFunc returns true if a and b contain the same items in the same order according to eq.


#### Example 
```go
{
	x := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	y := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	z := [][]byte{[]byte("a"), []byte("b"), []byte("d")}

	fmt.Println(xslices.EqualFunc(x, y, bytes.Equal))
	fmt.Println(xslices.EqualFunc(x[:2], y, bytes.Equal))
	fmt.Println(xslices.EqualFunc(z, y, bytes.Equal))

}
```

Output:
```text
true
false
false
```
<h3><a id="Fill"></a><samp>func <a href="#Fill">Fill</a>[T any](s []T, x T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L139">src</a></small></sub></h3>

Fill fills s with copies of x.


#### Example 
```go
{
	s := []int{1, 2, 3}
	xslices.Fill(s, 5)
	fmt.Println(s)

}
```

Output:
```text
[5 5 5]
```
<h3><a id="Filter"></a><samp>func <a href="#Filter">Filter</a>[T any](s []T, keep func(t T) bool) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L147">src</a></small></sub></h3>

Filter returns a slice containing only the elements of s for which keep() returns true in the
same order that they appeared in s.


#### Example 
```go
{
	s := []int{5, -9, -2, 1, -4, 8, 3}
	s = xslices.Filter(s, func(value int) bool {
		return value > 0
	})
	fmt.Println(s)

}
```

Output:
```text
[5 1 8 3]
```
<h3><a id="FilterInPlace"></a><samp>func <a href="#FilterInPlace">FilterInPlace</a>[T any](s []T, keep func(t T) bool) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L154">src</a></small></sub></h3>

FilterInPlace returns a slice containing only the elements of s for which keep() returns true in
the same order that they appeared in s. This is done in-place and so modifies the contents of s.
The modified slice is returned.


#### Example 
```go
{
	s := []int{5, -9, -2, 1, -4, 8, 3}
	s = xslices.FilterInPlace(s, func(value int) bool {
		return value > 0
	})
	fmt.Println(s)

}
```

Output:
```text
[5 1 8 3]
```
<h3><a id="Group"></a><samp>func <a href="#Group">Group</a>[T any, U comparable](s []T, f func(T) U) map[U][]T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L171">src</a></small></sub></h3>

Group returns a map from u to all items of s for which f(s[i]) returned u.


#### Example 
```go
{
	words := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}

	groups := xslices.Group(words, func(s string) rune {
		return ([]rune(s))[0]
	})

	for firstChar, group := range groups {
		fmt.Printf("%c: %v\n", firstChar, group)
	}

}
```

Unordered output:
```text
b: [bank beach]
g: [ghost goat group]
y: [yaw yew]
```
<h3><a id="Grow"></a><samp>func <a href="#Grow">Grow</a>[T any](s []T, n int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L183">src</a></small></sub></h3>

Grow grows s's capacity by reallocating, if necessary, to fit n more elements and returns the
modified slice. This does not change the length of s. After Grow(s, n), the following n
append()s to s will not need to reallocate.


#### Example 
```go
{
	s := make([]int, 0, 1)
	s = xslices.Grow(s, 4)
	fmt.Println(len(s))
	fmt.Println(cap(s))
	s = append(s, 1)
	addr := &s[0]
	s = append(s, 2)
	fmt.Println(addr == &s[0])
	s = append(s, 3)
	fmt.Println(addr == &s[0])
	s = append(s, 4)
	fmt.Println(addr == &s[0])

}
```

Output:
```text
0
4
true
true
true
```
<h3><a id="Index"></a><samp>func <a href="#Index">Index</a>[T comparable](s []T, x T) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L193">src</a></small></sub></h3>

Index returns the first index of x in s, or -1 if x is not in s.


#### Example 
```go
{
	s := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.Index(s, "b"))
	fmt.Println(xslices.Index(s, "c"))

}
```

Output:
```text
1
-1
```
<h3><a id="IndexFunc"></a><samp>func <a href="#IndexFunc">IndexFunc</a>[T any](s []T, f func(T) bool) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L204">src</a></small></sub></h3>

Index returns the first index in s for which f(s[i]) returns true, or -1 if there are no such
items.


#### Example 
```go
{
	s := []string{
		"blue",
		"green",
		"yellow",
		"gold",
		"red",
	}

	fmt.Println(xslices.IndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "g")
	}))
	fmt.Println(xslices.IndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "p")
	}))

}
```

Output:
```text
1
-1
```
<h3><a id="Insert"></a><samp>func <a href="#Insert">Insert</a>[T any](s []T, idx int, values ...T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L220">src</a></small></sub></h3>

Insert inserts the given values starting at index idx, shifting elements after idx to the right
and growing the slice to make room. Insert will expand the length of the slice up to its capacity
if it can, if this isn't desired then s should be resliced to have capacity equal to its length:

  s[:len(s):len(s)]

The time cost is O(n+m) where n is len(values) and m is len(s[idx:]).


#### Example 
```go
{
	s := []string{"a", "b", "c", "d", "e"}
	s = xslices.Insert(s, 3, "f", "g")
	fmt.Println(s)

}
```

Output:
```text
[a b c f g d e]
```
<h3><a id="Join"></a><samp>func <a href="#Join">Join</a>[T any](in ...[]T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L229">src</a></small></sub></h3>

Join joins together the contents of each in.


#### Example 
```go
{
	joined := xslices.Join(
		[]string{"a", "b", "c"},
		[]string{"x", "y"},
		[]string{"l", "m", "n", "o"},
	)

	fmt.Println(joined)

}
```

Output:
```text
[a b c x y l m n o]
```
<h3><a id="LastIndex"></a><samp>func <a href="#LastIndex">LastIndex</a>[T comparable](s []T, x T) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L242">src</a></small></sub></h3>

LastIndex returns the last index of x in s, or -1 if x is not in s.


#### Example 
```go
{
	s := []string{"a", "b", "a", "a", "b"}

	fmt.Println(xslices.LastIndex(s, "a"))
	fmt.Println(xslices.LastIndex(s, "c"))

}
```

Output:
```text
3
-1
```
<h3><a id="LastIndexFunc"></a><samp>func <a href="#LastIndexFunc">LastIndexFunc</a>[T any](s []T, f func(T) bool) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L253">src</a></small></sub></h3>

LastIndexFunc returns the last index in s for which f(s[i]) returns true, or -1 if there are no
such items.


#### Example 
```go
{
	s := []string{
		"blue",
		"green",
		"yellow",
		"gold",
		"red",
	}

	fmt.Println(xslices.LastIndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "g")
	}))
	fmt.Println(xslices.LastIndexFunc(s, func(s string) bool {
		return strings.HasPrefix(s, "p")
	}))

}
```

Output:
```text
3
-1
```
<h3><a id="Map"></a><samp>func <a href="#Map">Map</a>[T any, U any](s []T, f func(T) U) []U</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L263">src</a></small></sub></h3>

Map creates a new slice by applying f to each element of s.


#### Example 
```go
{
	toHalfFloat := func(x int) float32 {
		return float32(x) / 2
	}

	s := []int{1, 2, 3}
	floats := xslices.Map(s, toHalfFloat)
	fmt.Println(floats)

}
```

Output:
```text
[0.5 1 1.5]
```
<h3><a id="Partition"></a><samp>func <a href="#Partition">Partition</a>[T any](s []T, f func(t T) bool) int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L275">src</a></small></sub></h3>

Partition moves elements of s such that all elements for which f returns false are at the
beginning and all elements for which f returns true are at the end. It makes no other guarantees
about the final order of elements. Returns the index of the first element for which f returned
true, or len(s) if there wasn't one.


#### Example 
```go
{
	s := []int{11, 3, 4, 2, 7, 8, 0, 1, 14}

	xslices.Partition(s, func(x int) bool { return x%2 == 0 })

	fmt.Println(s)

}
```

Output:
```text
[11 3 1 7 2 8 0 4 14]
```
<h3><a id="Reduce"></a><samp>func <a href="#Reduce">Reduce</a>[T any, U any](s []T, initial U, f func(U, T) U) U</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L307">src</a></small></sub></h3>

Reduce reduces s to a single value using the reduction function f.


#### Example 
```go
{
	s := []int{3, 1, 2}

	sum := xslices.Reduce(s, 0, func(x, y int) int { return x + y })
	fmt.Println(sum)

	min := xslices.Reduce(s, math.MaxInt, xmath.Min[int])
	fmt.Println(min)

}
```

Output:
```text
6
1
```
<h3><a id="Remove"></a><samp>func <a href="#Remove">Remove</a>[T any](s []T, idx int, n int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L318">src</a></small></sub></h3>

Remove removes n elements from s starting at index idx and returns the modified slice. This
requires shifting the elements after the removed elements over, and so its cost is linear in the
number of elements shifted.


#### Example 
```go
{
	s := []int{1, 2, 3, 4, 5}
	s = xslices.Remove(s, 1, 2)
	fmt.Println(s)

}
```

Output:
```text
[1 4 5]
```
<h3><a id="RemoveUnordered"></a><samp>func <a href="#RemoveUnordered">RemoveUnordered</a>[T any](s []T, idx int, n int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L328">src</a></small></sub></h3>

RemoveUnordered removes n elements from s starting at index idx and returns the modified slice.
This is done by moving up to n elements from the end of the slice into the gap left by removal,
which is linear in n (rather than len(s)-idx as Remove() is), but does not preserve order of the
remaining elements.


#### Example 
```go
{
	s := []int{1, 2, 3, 4, 5}
	s = xslices.RemoveUnordered(s, 1, 1)
	fmt.Println(s)

	s = xslices.RemoveUnordered(s, 1, 2)
	fmt.Println(s)

}
```

Output:
```text
[1 5 3 4]
[1 4]
```
<h3><a id="Repeat"></a><samp>func <a href="#Repeat">Repeat</a>[T any](s T, n int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L340">src</a></small></sub></h3>

Repeat returns a slice with length n where every item is s.


#### Example 
```go
{
	s := xslices.Repeat("a", 4)
	fmt.Println(s)

}
```

Output:
```text
[a a a a]
```
<h3><a id="Reverse"></a><samp>func <a href="#Reverse">Reverse</a>[T any](s []T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L349">src</a></small></sub></h3>

Reverse reverses the elements of s in place.


#### Example 
```go
{
	s := []string{"a", "b", "c", "d", "e"}
	xslices.Reverse(s)
	fmt.Println(s)

}
```

Output:
```text
[e d c b a]
```
<h3><a id="Runs"></a><samp>func <a href="#Runs">Runs</a>[T any](s []T, same func(a, b T) bool) [][]T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L362">src</a></small></sub></h3>

Runs returns a slice of slices. The inner slices are contiguous runs of elements from s such that
same(a, b) returns true for any a and b in the run.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then same(a, c) must
also.

The returned slices use the same underlying array as s.


#### Example 
```go
{
	s := []int{2, 4, 0, 7, 1, 3, 9, 2, 8}

	parityRuns := xslices.Runs(s, func(a, b int) bool {
		return a%2 == b%2
	})

	fmt.Println(parityRuns)

}
```

Output:
```text
[[2 4 0] [7 1 3 9] [2 8]]
```
<h3><a id="Shrink"></a><samp>func <a href="#Shrink">Shrink</a>[T any](s []T, n int) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L382">src</a></small></sub></h3>

Shrink shrinks s's capacity by reallocating, if necessary, so that cap(s) <= len(s) + n.


#### Example 
```go
{
	s := make([]int, 3, 15)
	s[0] = 0
	s[1] = 1
	s[2] = 2

	fmt.Println(s)
	fmt.Println(cap(s))

	s = xslices.Shrink(s, 0)

	fmt.Println(s)
	fmt.Println(cap(s))

}
```

Output:
```text
[0 1 2]
15
[0 1 2]
3
```
<h3><a id="Unique"></a><samp>func <a href="#Unique">Unique</a>[T comparable](s []T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L396">src</a></small></sub></h3>

Unique returns a slice that contains only the first instance of each unique item in s, preserving
order.

Compact is more efficient if duplicates are already adjacent in s, for example if s is in sorted
order.


#### Example 
```go
{
	s := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := xslices.Unique(s)
	fmt.Println(unique)

}
```

Output:
```text
[a b c]
```
<h3><a id="UniqueInPlace"></a><samp>func <a href="#UniqueInPlace">UniqueInPlace</a>[T comparable](s []T) []T</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/xslices/xslices.go#L406">src</a></small></sub></h3>

UniqueInPlace returns a slice that contains only the first instance of each unique item in s,
preserving order. This is done in-place and so modifies the contents of s. The modified slice is
returned.

Compact is more efficient if duplicates are already adjacent in s, for example if s is in sorted
order.


#### Example 
```go
{
	s := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := xslices.UniqueInPlace(s)
	fmt.Println(unique)

}
```

Output:
```text
[a b c]
```
## Types

This section is empty.

