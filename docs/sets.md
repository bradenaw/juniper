# `package sets`

```
import "github.com/bradenaw/juniper/sets"
```

# Overview

package sets contains set operations like union, intersection, and difference.


# Index

<pre><a href="#Set">type Set</a></pre>
<pre>    <a href="#Difference">func Difference[T comparable](out, a, b Set[T]) Set[T]</a></pre>
<pre>    <a href="#Intersection">func Intersection[T comparable](out Set[T], sets ...Set[T]) Set[T]</a></pre>
<pre>    <a href="#Union">func Union[T any](out Set[T], sets ...Set[T]) Set[T]</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

# Types

## <a id="Set"></a><pre>type Set</pre>
```go
type Set[T any] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Len() int
	Iterate() iterator.Iterator[T]
}
```



<h2><a id="Difference"></a><pre>func Difference[T comparable](out, a, b <a href="#Set">Set</a>[T]) <a href="#Set">Set</a>[T]</pre></h2>

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
<h2><a id="Intersection"></a><pre>func Intersection[T comparable](out <a href="#Set">Set</a>[T], sets ...) <a href="#Set">Set</a>[T]</pre></h2>

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
<h2><a id="Union"></a><pre>func Union[T any](out <a href="#Set">Set</a>[T], sets ...) <a href="#Set">Set</a>[T]</pre></h2>

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
