# `package slices`

```
import "github.com/bradenaw/juniper/slices"
```

# Overview



# Index

<samp><a href="#All">func All[T any](x []T, f func(T) bool) bool</a></samp>
<samp><a href="#Any">func Any[T any](x []T, f func(T) bool) bool</a></samp>
<samp><a href="#Chunk">func Chunk[T any](x []T, chunkSize int) [][]T</a></samp>
<samp><a href="#Clear">func Clear[T any](x []T)</a></samp>
<samp><a href="#Clone">func Clone[T any](x []T) []T</a></samp>
<samp><a href="#Compact">func Compact[T comparable](x []T) []T</a></samp>
<samp><a href="#CompactFunc">func CompactFunc[T any](x []T, eq func(T, T) bool) []T</a></samp>
<samp><a href="#Count">func Count[T comparable](a []T, item T) int</a></samp>
<samp><a href="#CountFunc">func CountFunc[T any](a []T, f func(T) bool) int</a></samp>
<samp><a href="#Equal">func Equal[T comparable](a, b []T) bool</a></samp>
<samp><a href="#Fill">func Fill[T any](a []T, x T)</a></samp>
<samp><a href="#Filter">func Filter[T any](x []T, keep func(t T) bool) []T</a></samp>
<samp><a href="#Flatten">func Flatten[T any](x [][]T) []T</a></samp>
<samp><a href="#Grow">func Grow[T any](x []T, n int) []T</a></samp>
<samp><a href="#Index">func Index[T comparable](a []T, item T) int</a></samp>
<samp><a href="#IndexFunc">func IndexFunc[T any](a []T, f func(T) bool) int</a></samp>
<samp><a href="#Insert">func Insert[T any](x []T, idx int, values ...T) []T</a></samp>
<samp><a href="#Join">func Join[T any](in ...[]T) []T</a></samp>
<samp><a href="#LastIndex">func LastIndex[T comparable](a []T, item T) int</a></samp>
<samp><a href="#LastIndexFunc">func LastIndexFunc[T any](a []T, f func(T) bool) int</a></samp>
<samp><a href="#Map">func Map[T any, U any](x []T, f func(T) U) []U</a></samp>
<samp><a href="#Partition">func Partition[T any](x []T, f func(t T) bool)</a></samp>
<samp><a href="#Reduce">func Reduce[T any, U any](x []T, initial U, f func(U, T) U) U</a></samp>
<samp><a href="#Remove">func Remove[T any](x []T, idx int, n int) []T</a></samp>
<samp><a href="#Repeat">func Repeat[T any](x T, n int) []T</a></samp>
<samp><a href="#Reverse">func Reverse[T any](x []T)</a></samp>
<samp><a href="#Runs">func Runs[T any](x []T, same func(a, b T) bool) [][]T</a></samp>
<samp><a href="#Shrink">func Shrink[T any](x []T, n int) []T</a></samp>
<samp><a href="#Unique">func Unique[T comparable](x []T) []T</a></samp>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="All"></a><samp>func <a href="#All">All</a>[T any](x []T, f (T) bool) bool</samp></h2>

All returns true if f(x[i]) returns true for all i. Trivially, returns true if x is empty.


### Example 
```go
{
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	allOdd := slices.All([]int{1, 3, 5}, isOdd)
	fmt.Println(allOdd)

	allOdd = slices.All([]int{1, 3, 6}, isOdd)
	fmt.Println(allOdd)

}
```

Output:
```text
true
false
```
<h2><a id="Any"></a><samp>func <a href="#Any">Any</a>[T any](x []T, f (T) bool) bool</samp></h2>

Any returns true if f(x[i]) returns true for any i. Trivially, returns false if x is empty.


### Example 
```go
{
	isOdd := func(x int) bool {
		return x%2 != 0
	}

	anyOdd := slices.Any([]int{2, 3, 4}, isOdd)
	fmt.Println(anyOdd)

	anyOdd = slices.Any([]int{2, 4, 6}, isOdd)
	fmt.Println(anyOdd)

}
```

Output:
```text
true
false
```
<h2><a id="Chunk"></a><samp>func <a href="#Chunk">Chunk</a>[T any](x []T, chunkSize int) [][]T</samp></h2>

Chunk returns non-overlapping chunks of x. The last chunk will be smaller than chunkSize if
len(x) is not a multiple of chunkSize.

Returns an empty slice if len(x)==0. Panics if chunkSize <= 0.


### Example 
```go
{
	a := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	chunks := slices.Chunk(a, 3)
	fmt.Println(chunks)

}
```

Output:
```text
[[a b c] [d e f] [g h]]
```
<h2><a id="Clear"></a><samp>func <a href="#Clear">Clear</a>[T any](x []T)</samp></h2>

Clear fills x with the zero value of T.


### Example 
```go
{
	x := []int{1, 2, 3}
	slices.Clear(x)
	fmt.Println(x)

}
```

Output:
```text
[0 0 0]
```
<h2><a id="Clone"></a><samp>func <a href="#Clone">Clone</a>[T any](x []T) []T</samp></h2>

Clone creates a new slice and copies the elements of x into it.


### Example 
```go
{
	x := []int{1, 2, 3}
	cloned := slices.Clone(x)
	fmt.Println(cloned)

}
```

Output:
```text
[1 2 3]
```
<h2><a id="Compact"></a><samp>func <a href="#Compact">Compact</a>[T comparable](x []T) []T</samp></h2>

Compact removes adjacent duplicates from x in-place and returns the modified slice.


### Example 
```go
{
	x := []string{"a", "a", "b", "c", "c", "c", "a"}
	compacted := slices.Compact(x)
	fmt.Println(compacted)

}
```

Output:
```text
[a b c a]
```
<h2><a id="CompactFunc"></a><samp>func <a href="#CompactFunc">CompactFunc</a>[T any](x []T, eq (T, T) bool) []T</samp></h2>

CompactFunc removes adjacent duplicates from x in-place, preserving the first occurrence, using
the supplied eq function and returns the modified slice.


### Example 
```go
{
	x := []string{
		"bank",
		"beach",
		"ghost",
		"goat",
		"group",
		"yaw",
		"yew",
	}
	compacted := slices.CompactFunc(x, func(a, b string) bool {
		return a[0] == b[0]
	})
	fmt.Println(compacted)

}
```

Output:
```text
[bank ghost yaw]
```
<h2><a id="Count"></a><samp>func <a href="#Count">Count</a>[T comparable](a []T, item T) int</samp></h2>

Count returns the number of times item appears in a.


### Example 
```go
{
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(slices.Count(x, "a"))

}
```

Output:
```text
3
```
<h2><a id="CountFunc"></a><samp>func <a href="#CountFunc">CountFunc</a>[T any](a []T, f (T) bool) int</samp></h2>

Count returns the number of items in a for which f returns true.


<h2><a id="Equal"></a><samp>func <a href="#Equal">Equal</a>[T comparable](a, b []T) bool</samp></h2>

Equal returns true if a and b contain the same items in the same order.


### Example 
```go
{
	x := []string{"a", "b", "c"}
	y := []string{"a", "b", "c"}
	z := []string{"a", "b", "d"}

	fmt.Println(slices.Equal(x, y))
	fmt.Println(slices.Equal(x[:2], y))
	fmt.Println(slices.Equal(z, y))

}
```

Output:
```text
true
false
false
```
<h2><a id="Fill"></a><samp>func <a href="#Fill">Fill</a>[T any](a []T, x T)</samp></h2>

Fill fills a with copies of x.


### Example 
```go
{
	x := []int{1, 2, 3}
	slices.Fill(x, 5)
	fmt.Println(x)

}
```

Output:
```text
[5 5 5]
```
<h2><a id="Filter"></a><samp>func <a href="#Filter">Filter</a>[T any](x []T, keep (t T) bool) []T</samp></h2>

Filter filters the contents of x to only those for which keep() returns true. This is done
in-place and so modifies the contents of x. The modified slice is returned.


### Example 
```go
{
	x := []int{5, -9, -2, 1, -4, 8, 3}
	x = slices.Filter(x, func(value int) bool {
		return value > 0
	})
	fmt.Println(x)

}
```

Output:
```text
[5 1 8 3]
```
<h2><a id="Flatten"></a><samp>func <a href="#Flatten">Flatten</a>[T any](x [][]T) []T</samp></h2>

Flatten returns a slice containing all of the elements of all elements of x.


### Example 
```go
{
	x := [][]int{
		{0, 1, 2},
		{3, 4, 5, 6},
		{7},
	}

	fmt.Println(x)
	fmt.Println(slices.Flatten(x))

}
```

Output:
```text
[[0 1 2] [3 4 5 6] [7]]
[0 1 2 3 4 5 6 7]
```
<h2><a id="Grow"></a><samp>func <a href="#Grow">Grow</a>[T any](x []T, n int) []T</samp></h2>

Grow grows x's capacity by reallocating, if necessary, to fit n more elements and returns the
modified slice. This does not change the length of x. After Grow(x, n), the following n
append()s to x will not need to reallocate.


### Example 
```go
{
	x := make([]int, 0, 1)
	x = slices.Grow(x, 4)
	fmt.Println(len(x))
	fmt.Println(cap(x))
	x = append(x, 1)
	addr := &x[0]
	x = append(x, 2)
	fmt.Println(addr == &x[0])
	x = append(x, 3)
	fmt.Println(addr == &x[0])
	x = append(x, 4)
	fmt.Println(addr == &x[0])

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
<h2><a id="Index"></a><samp>func <a href="#Index">Index</a>[T comparable](a []T, item T) int</samp></h2>

Index returns the first index of item in a, or -1 if item is not in a.


### Example 
```go
{
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(slices.Index(x, "b"))
	fmt.Println(slices.Index(x, "c"))

}
```

Output:
```text
1
-1
```
<h2><a id="IndexFunc"></a><samp>func <a href="#IndexFunc">IndexFunc</a>[T any](a []T, f (T) bool) int</samp></h2>

Index returns the first index in a for which f(a[i]) returns true, or -1 if there are no such
items.


<h2><a id="Insert"></a><samp>func <a href="#Insert">Insert</a>[T any](x []T, idx int, values ...) []T</samp></h2>

Insert inserts the given values starting at index idx, shifting elements after idx to the right
and growing the slice to make room. Insert will expand the length of the slice up to its capacity
if it can, if this isn't desired then x should be resliced to have capacity equal to its length:

  x[:len(x):len(x)]

The time cost is O(n+m) where n is len(values) and m is len(x[idx:]).


### Example 
```go
{
	x := []string{"a", "b", "c", "d", "e"}
	x = slices.Insert(x, 3, "f", "g")
	fmt.Println(x)

}
```

Output:
```text
[a b c f g d e]
```
<h2><a id="Join"></a><samp>func <a href="#Join">Join</a>[T any](in ...) []T</samp></h2>

Join joins together the contents of each in.


### Example 
```go
{
	joined := slices.Join(
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
<h2><a id="LastIndex"></a><samp>func <a href="#LastIndex">LastIndex</a>[T comparable](a []T, item T) int</samp></h2>

LastIndex returns the last index of item in a, or -1 if item is not in a.


### Example 
```go
{
	x := []string{"a", "b", "a", "a", "b"}

	fmt.Println(slices.LastIndex(x, "a"))
	fmt.Println(slices.LastIndex(x, "c"))

}
```

Output:
```text
3
-1
```
<h2><a id="LastIndexFunc"></a><samp>func <a href="#LastIndexFunc">LastIndexFunc</a>[T any](a []T, f (T) bool) int</samp></h2>

LastIndexFunc returns the last index in a for which f(a[i]) returns true, or -1 if there are no
such items.


<h2><a id="Map"></a><samp>func <a href="#Map">Map</a>[T any, U any](x []T, f (T) U) []U</samp></h2>

Map creates a new slice by applying f to each element of x.


### Example 
```go
{
	toHalfFloat := func(x int) float32 {
		return float32(x) / 2
	}

	a := []int{1, 2, 3}
	floats := slices.Map(a, toHalfFloat)
	fmt.Println(floats)

}
```

Output:
```text
[0.5 1 1.5]
```
<h2><a id="Partition"></a><samp>func <a href="#Partition">Partition</a>[T any](x []T, f (t T) bool)</samp></h2>

Partition moves elements of x such that all elements for which f returns false are at the
beginning and all elements for which f returns true are at the end. It makes no other guarantees
about the final order of elements.


### Example 
```go
{
	a := []int{11, 3, 4, 2, 7, 8, 0, 1, 14}

	slices.Partition(a, func(x int) bool { return x%2 == 0 })

	fmt.Println(a)

}
```

Output:
```text
[11 3 1 7 2 8 0 4 14]
```
<h2><a id="Reduce"></a><samp>func <a href="#Reduce">Reduce</a>[T any, U any](x []T, initial U, f (U, T) U) U</samp></h2>

Reduce reduces x to a single value using the reduction function f.


### Example 
```go
{
	x := []int{3, 1, 2}

	sum := slices.Reduce(x, 0, func(x, y int) int { return x + y })
	fmt.Println(sum)

	min := slices.Reduce(x, math.MaxInt, xmath.Min[int])
	fmt.Println(min)

}
```

Output:
```text
6
1
```
<h2><a id="Remove"></a><samp>func <a href="#Remove">Remove</a>[T any](x []T, idx int, n int) []T</samp></h2>

Remove removes n elements from x starting at index idx and returns the modified slice. This
requires shifting the elements after the removed elements over, and so its cost is linear in the
number of elements shifted.


### Example 
```go
{
	x := []int{1, 2, 3, 4, 5}
	x = slices.Remove(x, 1, 2)
	fmt.Println(x)

}
```

Output:
```text
[1 4 5]
```
<h2><a id="Repeat"></a><samp>func <a href="#Repeat">Repeat</a>[T any](x T, n int) []T</samp></h2>

Repeat returns a slice with length n where every item is x.


### Example 
```go
{
	x := slices.Repeat("a", 4)
	fmt.Println(x)

}
```

Output:
```text
```
<h2><a id="Reverse"></a><samp>func <a href="#Reverse">Reverse</a>[T any](x []T)</samp></h2>

Reverse reverses the elements of x in place.


### Example 
```go
{
	x := []string{"a", "b", "c", "d", "e"}
	slices.Reverse(x)
	fmt.Println(x)

}
```

Output:
```text
[e d c b a]
```
<h2><a id="Runs"></a><samp>func <a href="#Runs">Runs</a>[T any](x []T, same (a, b T) bool) [][]T</samp></h2>

Runs returns a slice of slices. The inner slices are contiguous runs of elements from x such
that same(a, b) returns true for any a and b in the run.

same(a, a) must return true. If same(a, b) and same(b, c) both return true, then same(a, c) must
also.

The returned slices use the same underlying array as x.


### Example 
```go
{
	x := []int{2, 4, 0, 7, 1, 3, 9, 2, 8}

	parityRuns := slices.Runs(x, func(a, b int) bool {
		return a%2 == b%2
	})

	fmt.Println(parityRuns)

}
```

Output:
```text
[[2 4 0] [7 1 3 9] [2 8]]
```
<h2><a id="Shrink"></a><samp>func <a href="#Shrink">Shrink</a>[T any](x []T, n int) []T</samp></h2>

Shrink shrinks x's capacity by reallocating, if necessary, so that cap(x) <= len(x) + n.


<h2><a id="Unique"></a><samp>func <a href="#Unique">Unique</a>[T comparable](x []T) []T</samp></h2>

Unique removes duplicates from x in-place, preserving order, and returns the modified slice.

Compact is more efficient if duplicates are already adjacent in x, for example if x is in sorted
order.


### Example 
```go
{
	a := []string{"a", "b", "b", "c", "a", "b", "b", "c"}
	unique := slices.Unique(a)
	fmt.Println(unique)

}
```

Output:
```text
[a b c]
```
# Types

