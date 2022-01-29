# `package maps`

```
import "github.com/bradenaw/juniper/maps"
```

# Overview



# Index

<pre><a href="#Keys">func Keys[K comparable, V any](m map[K]V) []K</a></pre>
<pre><a href="#Values">func Values[K comparable, V any](m map[K]V) []V</a></pre>
<pre><a href="#Set">type Set</a></pre>
<pre>    <a href="#Add">func (s Set[T]) Add(item T)</a></pre>
<pre>    <a href="#Contains">func (s Set[T]) Contains(item T) bool</a></pre>
<pre>    <a href="#Iterate">func (s Set[T]) Iterate() iterator.Iterator[T]</a></pre>
<pre>    <a href="#Len">func (s Set[T]) Len() int</a></pre>
<pre>    <a href="#Remove">func (s Set[T]) Remove(item T)</a></pre>

# Constants

This section is empty.

# Variables

This section is empty.

# Functions

## <a id="Keys"></a><pre>func <a href="#Keys">Keys</a>[K comparable, V any](m map[K]V) []K</pre>

Keys returns the keys of m as a slice.


## <a id="Values"></a><pre>func <a href="#Values">Values</a>[K comparable, V any](m map[K]V) []V</pre>

Values returns the values of m as a slice.


# Types

## <a id="Set"></a><pre>type Set</pre>
```go
type Set[T comparable] map[T]struct{}
```

Set implements sets.Set for map[T]struct{}.


## <a id="Add"></a><pre>func (s <a href="#Set">Set</a>[T]) Add(item T)</pre>



## <a id="Contains"></a><pre>func (s <a href="#Set">Set</a>[T]) Contains(item T) bool</pre>



## <a id="Iterate"></a><pre>func (s <a href="#Set">Set</a>[T]) Iterate() <a href="./iterator.md#Iterator">iterator.Iterator</a>[T]</pre>



## <a id="Len"></a><pre>func (s <a href="#Set">Set</a>[T]) Len() int</pre>



## <a id="Remove"></a><pre>func (s <a href="#Set">Set</a>[T]) Remove(item T)</pre>



