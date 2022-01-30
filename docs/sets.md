# `package sets`

```
import "github.com/bradenaw/juniper/sets"
```

# Overview

package sets contains set operations like union, intersection, and difference.


# Index

<samp><a href="#Set">type Set</a></samp>

<samp>        <a href="#Difference">func Difference[T comparable](out, a, b Set[T]) Set[T]</a></samp>

<samp>        <a href="#Intersection">func Intersection[T comparable](out Set[T], sets ...Set[T]) Set[T]</a></samp>

<samp>        <a href="#Union">func Union[T any](out Set[T], sets ...Set[T]) Set[T]</a></samp>


# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

<h2><a id="Set"></a><samp>type Set</samp></h2>
```go
type Set[T any] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Len() int
	Iterate() iterator.Iterator[T]
}
```



<h2><a id="Difference"></a><samp>func Difference[T comparable](out, a, b <a href="#Set">Set</a>[T]) <a href="#Set">Set</a>[T]</samp></h2>

Difference adds to out all items that appear in a but not in b and returns out.


### Example 
```go
{
	a := maps.Set[int]{
		1: {},
		4: {},
		5: {},
	}
	b := maps.Set[int]{
		3: {},
		4: {},
	}

	out := make(maps.Set[int])

	difference := sets.Difference[int](out, a, b)

	fmt.Println(difference)

}
```

Output:
```text
map[1:{} 5:{}]
```
<h2><a id="Intersection"></a><samp>func Intersection[T comparable](out <a href="#Set">Set</a>[T], sets ...) <a href="#Set">Set</a>[T]</samp></h2>

Intersection adds to out all items that appear in all sets and returns out.


### Example 
```go
{
	a := maps.Set[int]{
		1: {},
		2: {},
		4: {},
	}
	b := maps.Set[int]{
		1: {},
		3: {},
		4: {},
	}
	c := maps.Set[int]{
		1: {},
		4: {},
		5: {},
	}

	out := make(maps.Set[int])

	intersection := sets.Intersection[int](out, a, b, c)

	fmt.Println(intersection)

}
```

Output:
```text
map[1:{} 4:{}]
```
<h2><a id="Union"></a><samp>func Union[T any](out <a href="#Set">Set</a>[T], sets ...) <a href="#Set">Set</a>[T]</samp></h2>

Union adds to out out all items from sets and returns out.


### Example 
```go
{
	a := maps.Set[int]{
		1: {},
		4: {},
	}
	b := maps.Set[int]{
		3: {},
		4: {},
	}
	c := maps.Set[int]{
		1: {},
		5: {},
	}

	out := make(maps.Set[int])

	union := sets.Union[int](out, a, b, c)

	fmt.Println(union)

}
```

Output:
```text
map[1:{} 3:{} 4:{} 5:{}]
```
