# `package maps`

```
import "github.com/bradenaw/juniper/maps"
```

# Overview



# Index

<samp><a href="#Keys">func Keys[K comparable, V any](m map[K]V) []K</a></samp>
<samp><a href="#Values">func Values[K comparable, V any](m map[K]V) []V</a></samp>
<samp><a href="#Set">type Set</a></samp>
<samp>    <a href="#Add">func (s Set[T]) Add(item T)</a></samp>
<samp>    <a href="#Contains">func (s Set[T]) Contains(item T) bool</a></samp>
<samp>    <a href="#Iterate">func (s Set[T]) Iterate() iterator.Iterator[T]</a></samp>
<samp>    <a href="#Len">func (s Set[T]) Len() int</a></samp>
<samp>    <a href="#Remove">func (s Set[T]) Remove(item T)</a></samp>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

<h2><a id="Keys"></a><samp>func <a href="#Keys">Keys</a>[K comparable, V any](m map[K]V) []K</samp></h2>

Keys returns the keys of m as a slice.


<h2><a id="Values"></a><samp>func <a href="#Values">Values</a>[K comparable, V any](m map[K]V) []V</samp></h2>

Values returns the values of m as a slice.


# Types

<h2><a id="Set"></a><samp>type Set</samp></h2>
```go
type Set[T comparable] map[T]struct{}
```

Set implements sets.Set for map[T]struct{}.


<h2><a id="Add"></a><samp>func (s <a href="#Set">Set</a>[T]) Add(item T)</samp></h2>



<h2><a id="Contains"></a><samp>func (s <a href="#Set">Set</a>[T]) Contains(item T) bool</samp></h2>



<h2><a id="Iterate"></a><samp>func (s <a href="#Set">Set</a>[T]) Iterate() <a href="./iterator.md#Iterator">iterator.Iterator</a>[T]</samp></h2>



<h2><a id="Len"></a><samp>func (s <a href="#Set">Set</a>[T]) Len() int</samp></h2>



<h2><a id="Remove"></a><samp>func (s <a href="#Set">Set</a>[T]) Remove(item T)</samp></h2>



