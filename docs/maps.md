# `package maps`

```
import "github.com/bradenaw/juniper/maps"
```

## Overview



## Index

<samp><a href="#Keys">func Keys[K comparable, V any](m map[K]V) []K</a></samp>

<samp><a href="#Values">func Values[K comparable, V any](m map[K]V) []V</a></samp>

<samp><a href="#Set">type Set</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Add">func (s Set[T]) Add(item T)</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Contains">func (s Set[T]) Contains(item T) bool</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Iterate">func (s Set[T]) Iterate() iterator.Iterator[T]</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Len">func (s Set[T]) Len() int</a></samp>

<samp>&nbsp;&nbsp;&nbsp;&nbsp;<a href="#Remove">func (s Set[T]) Remove(item T)</a></samp>


## Constants

This section is empty.

## Variables

This section is empty.

## Functions

<h3><a id="Keys"></a><samp>func <a href="#Keys">Keys</a>[K comparable, V any](m map[K]V) []K</samp></h3>

Keys returns the keys of m as a slice.


<h3><a id="Values"></a><samp>func <a href="#Values">Values</a>[K comparable, V any](m map[K]V) []V</samp></h3>

Values returns the values of m as a slice.


## Types

<h3><a id="Set"></a><samp>type Set</samp></h3>
```go
type Set[T comparable] map[T]struct{}
```

Set implements sets.Set for map[T]struct{}.


<h3><a id="Add"></a><samp>func (s <a href="#Set">Set</a>[T]) Add(item T)</samp></h3>



<h3><a id="Contains"></a><samp>func (s <a href="#Set">Set</a>[T]) Contains(item T) bool</samp></h3>



<h3><a id="Iterate"></a><samp>func (s <a href="#Set">Set</a>[T]) Iterate() <a href="./iterator.html#Iterator">iterator.Iterator</a>[T]</samp></h3>



<h3><a id="Len"></a><samp>func (s <a href="#Set">Set</a>[T]) Len() int</samp></h3>



<h3><a id="Remove"></a><samp>func (s <a href="#Set">Set</a>[T]) Remove(item T)</samp></h3>



