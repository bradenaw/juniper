# `package sets`

```
import "github.com/bradenaw/juniper/sets"
```

## Overview

Package sets contains set operations like union, intersection, and difference.


## Index

<samp><a href="#Intersects">func Intersects[T any](sets ...Set[T]) bool</a></samp>

<samp><a href="#Map">type Map</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Add">func (s Map[T]) Add(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (s Map[T]) Contains(item T) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (s Map[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#IterateInternal">func (s Map[T]) IterateInternal(f func(T) bool)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (s Map[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Remove">func (s Map[T]) Remove(item T)</a></samp>

<samp><a href="#Set">type Set</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Difference">func Difference[T any](out, a, b Set[T]) Set[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Intersection">func Intersection[T any](out Set[T], sets ...Set[T]) Set[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Union">func Union[T any](out Set[T], sets ...Set[T]) Set[T]</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="Intersects"></a><samp>func <a href="#Intersects">Intersects</a>[T any](sets ...<a href="#Set">Set</a>[T]) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L132">src</a></small></sub></h3>

Intersects returns true if the given sets have any items in common.


#### Example 
```go
{
	a := sets.Map[int]{
		1: {},
		2: {},
	}
	b := sets.Map[int]{
		1: {},
		3: {},
	}
	c := sets.Map[int]{
		3: {},
		4: {},
	}

	fmt.Println(sets.Intersects[int](a, b))
	fmt.Println(sets.Intersects[int](b, c))
	fmt.Println(sets.Intersects[int](a, c))

}
```

Output:
```text
true
true
false
```
## Types

<h3><a id="Map"></a><samp>type Map</samp></h3>
```go
type Map[T comparable] map[T]struct{}
```

Map implements sets.Set for map[T]struct{}.


<h3><a id="Add"></a><samp>func (s <a href="#Map">Map</a>[T]) Add(item T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L14">src</a></small></sub></h3>



<h3><a id="Contains"></a><samp>func (s <a href="#Map">Map</a>[T]) Contains(item T) bool</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L22">src</a></small></sub></h3>



<h3><a id="Iterate"></a><samp>func (s <a href="#Map">Map</a>[T]) Iterate() <a href="./iterator.html#Iterator">iterator.Iterator</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L53">src</a></small></sub></h3>



<h3><a id="IterateInternal"></a><samp>func (s <a href="#Map">Map</a>[T]) IterateInternal(f func(T) bool)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L57">src</a></small></sub></h3>



<h3><a id="Len"></a><samp>func (s <a href="#Map">Map</a>[T]) Len() int</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L27">src</a></small></sub></h3>



<h3><a id="Remove"></a><samp>func (s <a href="#Map">Map</a>[T]) Remove(item T)</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L18">src</a></small></sub></h3>



<h3><a id="Set"></a><samp>type Set</samp></h3>
```go
type Set[T any] interface {
	Add(item T)
	Remove(item T)
	Contains(item T) bool
	Len() int
	Iterate() iterator.Iterator[T]
}
```

Set is a minimal interface to a set. It is implemented by sets.Map and container/tree.Set, among
others.


<h3><a id="Difference"></a><samp>func Difference[T any](out, a, b <a href="#Set">Set</a>[T]) <a href="#Set">Set</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L155">src</a></small></sub></h3>

Difference adds to out all items that appear in a but not in b and returns out.


#### Example 
```go
{
	a := sets.Map[int]{
		1: {},
		4: {},
		5: {},
	}
	b := sets.Map[int]{
		3: {},
		4: {},
	}

	out := make(sets.Map[int])

	difference := sets.Difference[int](out, a, b)

	fmt.Println(difference)

}
```

Output:
```text
map[1:{} 5:{}]
```
<h3><a id="Intersection"></a><samp>func Intersection[T any](out <a href="#Set">Set</a>[T], sets ...<a href="#Set">Set</a>[T]) <a href="#Set">Set</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L110">src</a></small></sub></h3>

Intersection adds to out all items that appear in all sets and returns out.


#### Example 
```go
{
	a := sets.Map[int]{
		1: {},
		2: {},
		4: {},
	}
	b := sets.Map[int]{
		1: {},
		3: {},
		4: {},
	}
	c := sets.Map[int]{
		1: {},
		4: {},
		5: {},
	}

	out := make(sets.Map[int])

	intersection := sets.Intersection[int](out, a, b, c)

	fmt.Println(intersection)

}
```

Output:
```text
map[1:{} 4:{}]
```
<h3><a id="Union"></a><samp>func Union[T any](out <a href="#Set">Set</a>[T], sets ...<a href="#Set">Set</a>[T]) <a href="#Set">Set</a>[T]</samp><sub class="float-right"><small><a href="https://github.com/bradenaw/juniper/blob/main/sets/sets.go#L99">src</a></small></sub></h3>

Union adds to out out all items from sets and returns out.


#### Example 
```go
{
	a := sets.Map[int]{
		1: {},
		4: {},
	}
	b := sets.Map[int]{
		3: {},
		4: {},
	}
	c := sets.Map[int]{
		1: {},
		5: {},
	}

	out := make(sets.Map[int])

	union := sets.Union[int](out, a, b, c)

	fmt.Println(union)

}
```

Output:
```text
map[1:{} 3:{} 4:{} 5:{}]
```
